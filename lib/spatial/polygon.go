package spatial

// Polygon is a data type for storing simple polygons.
type Polygon []Line

func (p Polygon) String() string {
	return stringPolygon(p)
}

func (p Polygon) ClipToBBox(bbox BBox) []Geom {
	return p.clipToBBox(bbox)
}
