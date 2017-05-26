package spatial

// Polygon is a data type for storing simple polygons: One outer ring and an arbitrary number of inner rings.
type Polygon []Line

func (p Polygon) ClipToBBox(nw, se Point) []Geom {
	onw, ose := p[0].BBox()
	if onw.X() >= nw.X() && onw.Y() >= nw.Y() && ose.X() <= se.X() && ose.Y() <= se.Y() {
		// geometry fully inside
		geom, _ := NewGeom(p)
		return []Geom{geom}
	}

	var newGeoms []Geom
	// Cut outer ring
	for _, line := range p[0].ClipToBBox(nw, se) {
		ls, err := line.LineString()
		if err != nil {
			panic("clipping a line didn't return a line")
		}
		poly, err := NewGeom(Polygon{ls})
		if err != nil {
			panic("constructing a geometry with polygon failed")
		}
		newGeoms = append(newGeoms, poly)
	}
	return newGeoms
}
