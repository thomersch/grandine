package spatial

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type Point struct{ X, Y float64 }

func (p *Point) Project(proj ConvertFunc) {
	np := proj(*p)
	p.X = np.X
	p.Y = np.Y
}

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
	// TODO: Have a possibility to reduce precision to a limited number of decimals.
	b = strconv.AppendFloat(b, p.X, 'f', -1, 64)
	b = append(b, ',')
	b = strconv.AppendFloat(b, p.Y, 'f', -1, 64)
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

	var allsegs = make([]Segment, 0, len(poly)*10) // preallocating for 10 segments per ring
	for _, ln := range poly {
		allsegs = append(allsegs, ln.Segments()...)
		allsegs = append(allsegs, Segment{ln[len(ln)-1], ln[0]}) // closing segment
	}

	var (
		outTestPoint = Point{bbox.SW.X - 1, bbox.SW.Y - 1}
		l            = Line{p, outTestPoint}
	)

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

func (p Point) Copy() Projectable {
	return &Point{X: p.X, Y: p.Y}
}

const earthRadiusMeters = 6371000 // WGS84, optimized for minimal square relative error

// DistanceTo returns the distance in meters between the points.
func (p *Point) HaversineDistance(p2 *Point) float64 {
	lat1 := degToRad(p.Y)
	lon1 := degToRad(p.X)
	lat2 := degToRad(p2.Y)
	lon2 := degToRad(p2.X)

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c * earthRadiusMeters
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
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
