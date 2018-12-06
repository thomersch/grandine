package spatial

// Polygon is a data type for storing simple polygons.
type Polygon []Line

func (p Polygon) Project(proj ConvertFunc) {
	for ri := range p {
		for i := range p[ri] {
			p[ri][i] = proj(p[ri][i])
		}
	}
}

func (p Polygon) Copy() Projectable {
	var np Polygon
	for _, ring := range p {
		np = append(np, ring.Copy().(Line))
	}
	return np
}

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

func (p Polygon) FixWinding() {
	for n, ring := range p {
		if n == 0 {
			// First ring must be outer.
			if ring.Clockwise() {
				ring.Reverse()
			}
			continue
		}
		// Compare in how many rings the point is located.
		// If the number is odd, it's a hole.
		var inrings int
		for ninner, cring := range p {
			if n == ninner {
				continue
			}
			if ring[0].InPolygon(Polygon{cring}) {
				inrings++
			}
		}
		if (inrings%2 == 0 && ring.Clockwise()) || (inrings%2 == 1 && !ring.Clockwise()) {
			ring.Reverse()
		}
	}
}

func (p Polygon) ValidTopology() bool {
	return len(p.topologyErrors()) == 0
}

type segErr struct {
	Ring int
	Seg  int
}

func (p Polygon) topologyErrors() (errSegments []segErr) {
	for nRing, ring := range p {
		for nSeg, seg := range ring.SegmentsWithClosing() {
			for nSegCmp, segCmp := range ring.SegmentsWithClosing() {
				if nSeg == nSegCmp {
					continue
				}
				ipt, has := seg.Intersection(segCmp)
				if has && (ipt != seg[0] && ipt != seg[1]) {
					errSegments = append(errSegments, segErr{Ring: nRing, Seg: nSeg})
				}
			}
		}
	}
	return
}
