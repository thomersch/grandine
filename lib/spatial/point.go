package spatial

import "math"

type Point [2]float64

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) Y() float64 {
	return p[1]
}

func (p Point) InBBox(nw, se Point) bool {
	return nw[0] <= p[0] && se[0] >= p[0] &&
		nw[1] <= p[1] && se[1] >= p[1]
}

func (p Point) ClipToBBox(nw, se Point) []Geom {
	if p.InBBox(nw, se) {
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
	pbnw, pbse := poly[0].BBox()
	if !p.InBBox(pbnw, pbse) {
		return false
	}

	outTestPoint := Point{pbnw[0] - 1, pbnw[1] - 1}

	var allsegs []Segment
	for _, ln := range poly {
		allsegs = append(allsegs, ln.Segments()...)
	}

	l := Line{p, outTestPoint}
	if len(l.Intersections(allsegs))%2 == 0 {
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
