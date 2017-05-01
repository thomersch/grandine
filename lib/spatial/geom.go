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
	GeomTypeLineString          = 2
	GeomTypePolygon             = 3
	GeomTypeInvalid
)

type Geom struct {
	typ GeomType
	g   interface{}
}

func NewGeom(g interface{}) (Geom, error) {
	switch g.(type) {
	case Point:
		return Geom{typ: GeomTypePoint, g: g}, nil
	case []Point:
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
		var ls []Point
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
	err := binary.Read(r, endianness, &wkbEndianness)
	if err != nil {
		return err
	}
	if wkbEndianness != 1 {
		return errors.New("only little endian is supported")
	}

	var (
		gt  GeomType
		npg interface{}
		ng  Geom
	)
	err = binary.Read(r, endianness, &gt)
	if err != nil {
		return err
	}
	switch gt {
	case GeomTypePoint:
		npg, err = wkbReadPoint(r)
	case GeomTypeLineString:
		npg, err = wkbReadLineString(r)
	default:
		panic("not implemented yet")
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
	var buf bytes.Buffer
	binary.Write(&buf, endianness, uint8(1)) // little endian
	binary.Write(&buf, endianness, g.Typ())  // geometry type

	switch g.Typ() {
	case GeomTypePoint:
		p, err := g.Point()
		if err != nil {
			return nil, err
		}
		wkbWritePoint(&buf, p)
	case GeomTypeLineString:
		ls, err := g.LineString()
		if err != nil {
			return nil, err
		}
		wkbWriteLineString(&buf, ls)
	case GeomTypePolygon:
		poly, err := g.Polygon()
		if err != nil {
			return nil, err
		}
		wkbWritePolygon(&buf, poly)
	default:
		panic("not implemented yet")
	}
	return buf.Bytes(), nil
}

func (g *Geom) Typ() GeomType {
	return g.typ
}

func (g *Geom) Point() (Point, error) {
	geom, ok := g.g.(Point)
	if !ok {
		return Point{}, errors.New("geometry is not a Point")
	}
	return geom, nil
}

func (g *Geom) LineString() ([]Point, error) {
	geom, ok := g.g.([]Point)
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

type Feature struct {
	Type     string                 `json:"type"`
	Props    map[string]interface{} `json:"properties"`
	Geometry Geom                   `json:"geometry"`
}

func (f *Feature) Properties() map[string]interface{} {
	return f.Props
}

type FeatureCollection struct {
	Features []Feature `json:"features"`
}

func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	wfc := struct {
		Type     string    `json:"type"`
		Features []Feature `json:"features"`
	}{
		Type:     "FeatureCollection",
		Features: fc.Features,
	}
	return json.Marshal(wfc)
}
