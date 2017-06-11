package mvt

import (
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

// Some of the contents are independent of the MVT implementation and could probably be moved to
// a separate package.

type TileID struct {
	Z, X, Y int
}

func (t TileID) BBox() (NW, SE spatial.Point) {
	return t.NW(), TileID{X: t.X + 1, Y: t.Y + 1, Z: t.Z}.NW()
}

func (t TileID) NW() spatial.Point {
	n := math.Pow(2, float64(t.Z))
	lonDeg := float64(t.X)/n*360 - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(t.Y)/n)))
	latDeg := latRad * 180 / math.Pi
	return spatial.Point{lonDeg, latDeg}
}

func TileName(p spatial.Point, zoom int) TileID {
	var (
		zf     = float64(zoom)
		latDeg = float64(p.Y() * math.Pi / 180)
	)
	return TileID{
		X: int(math.Floor((float64(p.X()) + 180) / 360 * math.Exp2(zf))),
		Y: int(math.Floor((1 - math.Log(math.Tan(latDeg)+1/math.Cos(latDeg))/math.Pi) / 2 * math.Exp2(zf))),
		Z: zoom,
	}
}
