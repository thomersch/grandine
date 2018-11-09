package spatial

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type Point struct{ X, Y float64 }

func (p Point) InBBox(b BBox) bool {
	return b.SW.X <= p.X && b.NE.X >= p.X &&
		b.SW.Y <= p.Y && b.NE.Y >= p.Y
}

func (p Point) ClipToBBox(b BBox) []Geom {
	if p.InBBox(b) {
		g := MustNewGeom(p)
		return []Geom{g}
	}
	return []Geom{}
}

func (p Point) String() string {
	return fmt.Sprintf("(X: %v, Y: %v)", p.X, p.Y)
}

func (p Point) MarshalJSON() ([]byte, error) {
	var b = make([]byte, 1, 30)
	b[0] = '['
	b = strconv.AppendFloat(b, p.X, 'f', 8, 64)
	b = append(b, ',')
	b = strconv.AppendFloat(b, p.Y, 'f', 8, 64)
	return append(b, ']'), nil
}

func (p *Point) SetX(x float64) {
	p.X = x
}

func (p *Point) SetY(y float64) {
	p.Y = y
}

func (p *Point) UnmarshalJSON(buf []byte) error {
	var arrayPt [2]float64
	if err := json.Unmarshal(buf, &arrayPt); err != nil {
		return err
	}
	p.X = arrayPt[0]
	p.Y = arrayPt[1]
	return nil
}

const pointPrecision = 8

func (p Point) RoundedCoords() Point {
	return Point{
		roundWithPrecision(p.X, pointPrecision),
		roundWithPrecision(p.Y, pointPrecision),
	}
}

func (p Point) InPolygon(poly Polygon) bool {
	bbox := poly[0].BBox()
	if !p.InBBox(bbox) {
		return false
	}

	outTestPoint := Point{bbox.SW.X - 1, bbox.SW.Y - 1}

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
