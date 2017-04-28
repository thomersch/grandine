package spatial

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

var endianness = binary.LittleEndian

type FeatureCollection struct {
	Features []Feature
}

type Feature struct {
	// For now this is a GeoJSON feature
	Type     string
	Props    map[string]interface{} `json:"properties"`
	Geometry struct {
		Type        string
		Coordinates Coords
	}
}

type Point [2]float64

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) Y() float64 {
	return p[1]
}

type Coords struct {
	b []byte
}

func (c *Coords) UnmarshalJSON(buf []byte) error {
	c.b = buf
	return nil
}

func (c *Coords) Point() (Point, error) {
	var p Point
	err := json.Unmarshal(c.b, &p)
	return p, err
}

func (c *Coords) LineString() ([]Point, error) {
	var ls []Point
	err := json.Unmarshal(c.b, &ls)
	return ls, err
}

func (c *Coords) Polygon() ([][]Point, error) {
	var poly [][]Point
	err := json.Unmarshal(c.b, &poly)
	return poly, err
}

type PropertyRetriever interface {
	Properties() map[string]interface{}
}

func (f *Feature) MarshalWKB() ([]byte, error) {
	if endianness != binary.LittleEndian {
		return nil, errors.New("only little endian is supported")
	}
	var buf bytes.Buffer
	binary.Write(&buf, endianness, uint8(1)) // little endian
	binary.Write(&buf, endianness, f.Typ())  // geometry type

	switch f.Typ() {
	case GeomTypePoint:
		p, err := f.Geometry.Coordinates.Point()
		if err != nil {
			return nil, err
		}
		binary.Write(&buf, endianness, p.X)
		binary.Write(&buf, endianness, p.Y)
	}
	return buf.Bytes(), nil
}

type GeomType uint32

const (
	GeomTypePoint      GeomType = 1
	GeomTypeLineString          = 2
	GeomTypePolygon             = 3
	GeomTypeInvalid
)

func (f *Feature) Typ() GeomType {
	switch f.Geometry.Type {
	case "Point":
		return GeomTypePoint
	case "LineString":
		return GeomTypeLineString
	case "Polygon":
		return GeomTypePolygon
	}
	return GeomTypeInvalid
}

func (f *Feature) Properties() map[string]interface{} {
	return f.Props
}
