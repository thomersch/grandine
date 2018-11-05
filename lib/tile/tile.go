package tile

import (
	"fmt"
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

const (
	wgs84LatMax = 85.0511287
	wgs84LonMax = 180
)

type ID struct {
	X, Y, Z int
}

func (t ID) BBox() spatial.BBox {
	nw := t.NW()
	se := ID{X: t.X + 1, Y: t.Y + 1, Z: t.Z}.NW()
	return spatial.BBox{SW: spatial.Point{nw.X, se.Y}, NE: spatial.Point{se.X, nw.Y}}
}

func (t ID) NW() spatial.Point {
	n := math.Pow(2, float64(t.Z))
	lonDeg := float64(t.X)/n*360 - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(t.Y)/n)))
	latDeg := latRad * 180 / math.Pi
	return spatial.Point{lonDeg, latDeg}
}

func (t ID) String() string {
	return fmt.Sprintf("%v/%v/%v", t.Z, t.X, t.Y)
}

func TileName(p spatial.Point, zoom int) ID {
	var (
		zf     = float64(zoom)
		latDeg = float64(floatBetween(-1*wgs84LatMax, wgs84LatMax, p.Y) * math.Pi / 180)
	)
	return ID{
		X: between(0, int(math.Exp2(zf)-1),
			int(math.Floor((float64(p.X)+180)/360*math.Exp2(zf)))),
		Y: between(0, int(math.Exp2(zf)-1),
			int(math.Floor((1-math.Log(math.Tan(latDeg)+1/math.Cos(latDeg))/math.Pi)/2*math.Exp2(zf)))),
		Z: zoom,
	}
}

// Resolution determines the minimal describable value inside a tile.
func Resolution(zoomlevel, extent int) float64 {
	return 360 / (math.Pow(2, float64(zoomlevel)) * float64(extent))
}

func between(min, max, v int) int {
	return int(math.Max(math.Min(float64(v), float64(max)), float64(min)))
}

func floatBetween(min, max, v float64) float64 {
	return math.Max(math.Min(v, max), min)
}
