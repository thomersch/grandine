package spatial

// Polygon is a data type for storing simple polygons: One outer ring and an arbitrary number of inner rings.
type Polygon []Line

func (p Polygon) ClipToBBox(nw, se Point) []Geom {
	// TODO: inner ring handling

	// Is outer ring fully inside?
	oNW, oSE := p[0].BBox()
	if oNW.X() >= nw.X() && oNW.Y() >= nw.Y() && oSE.X() <= se.X() && oSE.Y() <= se.Y() {
		geom, _ := NewGeom(p)
		return []Geom{geom}
	}

	var (
		newGeoms []Geom
		lns      []Line
	)
	// Cut outer ring
	gms := p[0].ClipToBBox(nw, se)
	for _, gm := range gms {
		lns = append(lns, gm.g.(Line))
	}
	polys := PolygonsFromLines(lns)
	for _, poly := range polys {
		ng, err := NewGeom(poly)
		if err != nil {
			panic("constructing a geometry with polygon failed")
		}
		newGeoms = append(newGeoms, ng)
	}
	return newGeoms
}

func PolygonsFromLines(ls []Line) []Polygon {
	var mlines []Line
	for lni := range ls {
		var merged bool
		for mi := range mlines {
			if mlines[mi].IsExtendedBy(ls[lni]) {
				mlines[mi] = MergeLines(mlines[mi], ls[lni])
				merged = true
			}
		}
		if !merged {
			mlines = append(mlines, ls[lni])
		}
	}
	var polys []Polygon
	for i := range mlines {
		if !mlines[i].Closed() {
			mlines[i] = append(mlines[i], mlines[i][0])
		}
		polys = append(polys, Polygon{mlines[i]})
	}
	return polys
}
