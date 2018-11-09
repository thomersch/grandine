package spatial

// Polygon is a data type for storing simple polygons.
type Polygon []Line

func (p Polygon) String() string {
	return p.string()
}

func (p Polygon) ClipToBBox(bbox BBox) []Geom {
	return p.clipToBBox(bbox)
}

func (p Polygon) Rewind() {
	for _, ring := range p {
		ring.Reverse()
	}
}
