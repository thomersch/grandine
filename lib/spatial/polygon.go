package spatial

import "fmt"

// Polygon is a data type for storing simple polygons.
type Polygon []Line

func (p Polygon) String() string {
	s := "Polygon{"
	for _, line := range p {
		s += fmt.Sprintf("%v", line)
	}
	return s + "}"
}

func (p Polygon) ClipToBBox(bbox BBox) []Geom {
	return p.clipToBBox(bbox)
}
