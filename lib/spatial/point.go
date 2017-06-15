package spatial

import "math"

type Point [2]float64

func (p Point) X() float64 {
	return p[0]
}

func (p Point) Y() float64 {
	return p[1]
}

func (p Point) InBBox(sw, ne Point) bool {
	return sw[0] <= p[0] && ne[0] >= p[0] &&
		sw[1] <= p[1] && ne[1] >= p[1]
}

func (p Point) ClipToBBox(sw, ne Point) []Geom {
	if p.InBBox(sw, ne) {
		g, _ := NewGeom(p)
		return []Geom{g}
	}
	return []Geom{}
}

const pointPrecision = 8

func (p Point) RoundedCoords() Point {
	return Point{
		roundWithPrecision(p[0], pointPrecision),
		roundWithPrecision(p[1], pointPrecision),
	}
}

func (p Point) InPolygon(poly Polygon) bool {
	pbsw, pbne := poly[0].BBox()
	if !p.InBBox(pbsw, pbne) {
		return false
	}

	outTestPoint := Point{pbsw[0] - 1, pbsw[1] - 1}

	var allsegs []Segment
	for _, ln := range poly {
		allsegs = append(allsegs, ln.Segments()...)
	}

	l := Line{p, outTestPoint}
	intersections := l.Intersections(allsegs)
	if len(intersections)%2 == 0 {
		for _, itsct := range intersections {
			// if the intersection is exactly the tested point, the point is actually inside
			// TODO: test more
			if itsct == p {
				return true
			}
		}
		return false
	}
	return true
}

func round(v float64) float64 {
	if v < 0 {
		return math.Ceil(v - 0.5)
	}
	return math.Floor(v + 0.5)
}

func roundWithPrecision(v float64, decimals int) float64 {
	s := math.Pow(10, float64(decimals))
	return round(v*s) / s
}
