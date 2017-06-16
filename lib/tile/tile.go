package tile

import (
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

type ID struct {
	X, Y, Z int
}

func (t ID) BBox() spatial.BBox {
	nw := t.NW()
	se := ID{X: t.X + 1, Y: t.Y + 1, Z: t.Z}.NW()
	return spatial.BBox{spatial.Point{nw.X(), se.Y()}, spatial.Point{se.X(), nw.Y()}}
}

func (t ID) NW() spatial.Point {
	n := math.Pow(2, float64(t.Z))
	lonDeg := float64(t.X)/n*360 - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(t.Y)/n)))
	latDeg := latRad * 180 / math.Pi
	return spatial.Point{lonDeg, latDeg}
}

func TileName(p spatial.Point, zoom int) ID {
	var (
		zf     = float64(zoom)
		latDeg = float64(p.Y() * math.Pi / 180)
	)
	return ID{
		X: int(math.Floor((float64(p.X()) + 180) / 360 * math.Exp2(zf))),
		Y: int(math.Floor((1 - math.Log(math.Tan(latDeg)+1/math.Cos(latDeg))/math.Pi) / 2 * math.Exp2(zf))),
		Z: zoom,
	}
}
