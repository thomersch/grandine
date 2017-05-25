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

type Point [2]float64

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) Y() float64 {
	return p[1]
}

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

func NewGeom(g interface{}) (Geom, error) {
	switch geom := g.(type) {
	case Point:
		return Geom{typ: GeomTypePoint, g: g}, nil
	case []Point:
		return Geom{typ: GeomTypeLineString, g: Line(geom)}, nil
	case Line:
		return Geom{typ: GeomTypeLineString, g: g}, nil
	case [][]Point:
		return Geom{typ: GeomTypePolygon, g: g}, nil
	default:
		return Geom{}, fmt.Errorf("unknown input geom type: %T", g)
	}
}

type geoJSONGeom struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
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
		var poly [][]Point
		if err = json.Unmarshal(wg.Coordinates, &poly); err != nil {
			return err
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
		var ls []Point
		ls, err = g.LineString()
		if err != nil {
			return nil, err
		}
		err = wkbWriteLineString(&buf, ls)
	case GeomTypePolygon:
		var poly [][]Point
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

func (g *Geom) Polygon() ([][]Point, error) {
	geom, ok := g.g.([][]Point)
	if !ok {
		return nil, errors.New("geometry is not a Polygon")
	}
	return geom, nil
}

func (g *Geom) BBox() (nw, se Point) {
	switch gm := g.g.(type) {
	case Point:
		return gm, gm
	case Line:
		return ringBBox(gm)
	case [][]Point:
		var bboxPoints []Point
		for _, ring := range gm {
			neb, seb := ringBBox(ring)
			bboxPoints = append(bboxPoints, neb, seb)
		}
		return ringBBox(bboxPoints)
	default:
		panic("unimplemented type")
	}
	return
}

// Clips a geometry and returns a cropped copy. Returns a slice, because clip might result in multiple sub-Geoms.
func (g *Geom) ClipToBBox(nw, se Point) []Geom {
	switch gm := g.g.(type) {
	case Point:
		if nw[0] <= gm[0] && se[0] >= gm[0] &&
			nw[1] <= gm[1] && se[1] >= gm[1] {
			return []Geom{*g}
		}
		return []Geom{}

	case Line:
		lsNW, lsSE := g.BBox()
		// Is linestring completely inside bbox?
		if nw[0] <= lsNW[0] && se[0] >= lsSE[0] &&
			nw[1] <= lsNW[1] && se[1] >= lsSE[1] {
			// no clipping necessary
			return []Geom{*g}
		}

		// Is linestring fully outside the bbox?
		if lsSE[0] < nw[0] || lsSE[1] < nw[1] || lsNW[0] > se[0] || lsNW[1] > se[1] {
			return []Geom{}
		}

		var cutsegs []Segment
		for _, seg := range gm.Segments() {
			if seg.FullyInBBox(nw, se) {
				cutsegs = append(cutsegs, seg)
				continue
			}
			for _, bbl := range BBoxBorders(nw, se) {
				if ipt, intersects := seg.Intersection(bbl); intersects {
					s1, s2 := seg.SplitAt(ipt)

					if s1.FullyInBBox(nw, se) && s1.Length() != 0 {
						cutsegs = append(cutsegs, s1)
						break
					}
					if s2.FullyInBBox(nw, se) && s2.Length() != 0 {
						cutsegs = append(cutsegs, s2)
						break
					}
				}
			}
		}
		var gms []Geom
		for _, ln := range NewLinesFromSegments(cutsegs) {
			gms = append(gms, Geom{typ: GeomTypeLineString, g: ln})
		}
		return gms
	default:
		panic("unknown geom type")
	}
	panic("falling through")
}
