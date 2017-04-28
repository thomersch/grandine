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

type Coords struct {
	b []byte
}

func (c *Coords) UnmarshalJSON(buf []byte) error {
	c.b = buf
	return nil
}

func (c *Coords) Point() [2]float64 {
	p := [2]float64{}
	err := json.Unmarshal(c.b, &p)
	if err != nil {
		panic(err)
	}
	return p
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
	case Point:
		p := f.Geometry.Coordinates.Point()
		binary.Write(&buf, endianness, p[0])
		binary.Write(&buf, endianness, p[1])
	}
	return buf.Bytes(), nil
}

type GeomType uint32

const (
	Point      GeomType = 1
	LineString          = 2
	Polygon             = 3
	Invalid
)

func (f *Feature) Typ() GeomType {
	switch f.Geometry.Type {
	case "Point":
		return Point
	case "LineString":
		return LineString
	case "Polygon":
		return Polygon
	}
	return Invalid
}

func (f *Feature) Properties() map[string]interface{} {
	return f.Props
}
