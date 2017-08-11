package spatial

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var endianness = binary.LittleEndian

type GeomType uint32

const (
	GeomTypePoint      GeomType = 1
	GeomTypeLineString GeomType = 2
	GeomTypePolygon    GeomType = 3
	GeomTypeInvalid
)

type Geom struct {
	typ GeomType
	g   interface{}
}

func MustNewGeom(g interface{}) Geom {
	ng, err := NewGeom(g)
	if err != nil {
		panic(err)
	}
	return ng
}

func NewGeom(g interface{}) (Geom, error) {
	switch geom := g.(type) {
	// Point
	case Point:
		return Geom{typ: GeomTypePoint, g: g}, nil

	// Line String
	case []Point:
		return Geom{typ: GeomTypeLineString, g: Line(geom)}, nil
	case Line:
		return Geom{typ: GeomTypeLineString, g: g}, nil

	// Polygon
	case [][]Point:
		var poly Polygon
		for _, ln := range geom {
			poly = append(poly, Line(ln))
		}
		return Geom{typ: GeomTypePolygon, g: poly}, nil
	case []Line:
		return Geom{typ: GeomTypePolygon, g: Polygon(geom)}, nil
	case Polygon:
		return Geom{typ: GeomTypePolygon, g: g}, nil
	default:
		return Geom{}, fmt.Errorf("unknown input geom type: %T", g)
	}
}

type geoJSONGeom struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

func (g Geom) String() string {
	// TODO: this could probably be replaced with a type assertion and direct call to g.g.String()
	return fmt.Sprintf("%v", g.g)
}

func (g *Geom) UnmarshalJSON(buf []byte) error {
	var wg geoJSONGeom
	err := json.Unmarshal(buf, &wg)
	if err != nil {
		return err
	}

	switch strings.ToLower(wg.Type) {
	case "point":
		g.typ = GeomTypePoint
		var p Point
		if err = json.Unmarshal(wg.Coordinates, &p); err != nil {
			return err
		}
		g.g = p
	case "linestring":
		g.typ = GeomTypeLineString
		var ls Line
		if err = json.Unmarshal(wg.Coordinates, &ls); err != nil {
			return err
		}
		g.g = ls
	case "polygon":
		g.typ = GeomTypePolygon
		var poly Polygon
		if err = json.Unmarshal(wg.Coordinates, &poly); err != nil {
			return err
		}
		// NOTE: the geojson decoder assumes OGC/RFC 7946 compliant winding order
		for pos := range poly {
			// remove last element from every ring as it is unnecessary
			poly[pos] = poly[pos][:len(poly[pos])-1]
		}
		g.g = poly
	default:
		return fmt.Errorf("unsupported geometry type: %s", wg.Type)
	}
	return nil
}

func (g Geom) MarshalJSON() ([]byte, error) {
	var wg geoJSONGeom

	switch g.typ {
	case GeomTypePoint:
		wg.Type = "Point"
		pt, err := g.Point()
		if err != nil {
			return nil, err
		}
		wg.Coordinates, err = json.Marshal(pt)
		if err != nil {
			return nil, err
		}
	case GeomTypeLineString:
		wg.Type = "LineString"
		ls, err := g.LineString()
		if err != nil {
			return nil, err
		}
		wg.Coordinates, err = json.Marshal(ls)
		if err != nil {
			return nil, err
		}
	case GeomTypePolygon:
		wg.Type = "Polygon"
		// TODO: each ring must end with first point in ring
		poly, err := g.Polygon()
		if err != nil {
			return nil, err
		}
		wg.Coordinates, err = json.Marshal(poly)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(&wg)
}

func (g *Geom) UnmarshalWKB(r io.Reader) error {
	var wkbEndianness uint8
	//TODO: read the byte directly
	err := binary.Read(r, endianness, &wkbEndianness)
	if err != nil {
		return err
	}
	if wkbEndianness != 1 {
		return errors.New("only little endian is supported")
	}

	var (
		npg interface{}
		ng  Geom
	)
	gt, err := wkbReadHeader(r)
	if err != nil {
		return err
	}
	switch gt {
	case GeomTypePoint:
		npg, err = wkbReadPoint(r)
	case GeomTypeLineString:
		npg, err = wkbReadLineString(r)
	case GeomTypePolygon:
		npg, err = wkbReadPolygon(r)
	default:
		return fmt.Errorf("unsupported GeomType: %v", gt)
	}
	if err != nil {
		return err
	}
	ng, err = NewGeom(npg)
	if err != nil {
		return err
	}
	g.typ = ng.typ
	g.g = ng.g
	return nil
}

// TODO: maybe MarshalWKB could take an io.Writer instead of returning a buffer?
func (g Geom) MarshalWKB() ([]byte, error) {
	if endianness != binary.LittleEndian {
		return nil, errors.New("only little endian is supported")
	}
	var (
		buf     bytes.Buffer
		typeBuf = make([]byte, 4)
	)
	_, err := buf.Write([]byte{1}) // little endian
	if err != nil {
		return nil, err
	}
	endianness.PutUint32(typeBuf, uint32(g.Typ()))
	_, err = buf.Write(typeBuf)
	if err != nil {
		return nil, err
	}

	switch g.Typ() {
	case GeomTypePoint:
		var p Point
		p, err = g.Point()
		if err != nil {
			return nil, err
		}
		err = wkbWritePoint(&buf, p)
	case GeomTypeLineString:
		var ls Line
		ls, err = g.LineString()
		if err != nil {
			return nil, err
		}
		err = wkbWriteLineString(&buf, ls)
	case GeomTypePolygon:
		var poly Polygon
		poly, err = g.Polygon()
		if err != nil {
			return nil, err
		}
		err = wkbWritePolygon(&buf, poly)
	default:
		return nil, fmt.Errorf("unsupported GeomType: %v", g.Typ())
	}
	return buf.Bytes(), err
}

// Typ returns the geometries type.
func (g *Geom) Typ() GeomType {
	return g.typ
}

// Point retruns the geometry as a point (check type with Typ()) first.
func (g *Geom) Point() (Point, error) {
	geom, ok := g.g.(Point)
	if !ok {
		return Point{}, errors.New("geometry is not a Point")
	}
	return geom, nil
}

func (g *Geom) LineString() (Line, error) {
	geom, ok := g.g.(Line)
	if !ok {
		return nil, errors.New("geometry is not a LineString")
	}
	return geom, nil
}

func (g *Geom) Polygon() (Polygon, error) {
	geom, ok := g.g.(Polygon)
	if !ok {
		return nil, errors.New("geometry is not a Polygon")
	}
	return geom, nil
}

func (g *Geom) BBox() BBox {
	switch gm := g.g.(type) {
	case Point:
		return BBox{gm, gm}
	case Line:
		return gm.BBox()
	case Polygon:
		var bboxPoints Line
		for _, ring := range gm {
			bb := ring.BBox()
			bboxPoints = append(bboxPoints, bb.SW, bb.NE)
		}
		return bboxPoints.BBox()
	default:
		panic("unimplemented type")
	}
}

func (g *Geom) In(bbox BBox) bool {
	return g.BBox().In(bbox)
}

type simplifiable interface {
	Simplify(e float64) interface{}
}

func (g *Geom) Simplify(e float64) Geom {
	switch gm := g.g.(type) {
	case Line:
		return Geom{typ: g.typ, g: gm.Simplify(e)}
	}
	return *g
}

type Clippable interface {
	// TODO: consider returning primitive geom, instead of Geom
	ClipToBBox(BBox) []Geom
}

// Clips a geometry and returns a cropped copy. Returns a slice, because clip might result in multiple sub-Geoms.
func (g *Geom) ClipToBBox(bbox BBox) []Geom {
	if gm, ok := g.g.(Clippable); ok {
		return gm.ClipToBBox(bbox)
	}
	panic("internal geometry needs to fulfill Clippable interface")
}
