package spatial

import "math"

type Point [2]float64

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) Y() float64 {
	return p[1]
}

func (p Point) ClipToBBox(nw, se Point) []Geom {
	if nw[0] <= p[0] && se[0] >= p[0] &&
		nw[1] <= p[1] && se[1] >= p[1] {
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
