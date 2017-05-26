package spatial

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
