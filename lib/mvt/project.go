package mvt

import (
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

type bbox struct {
	NW, SE point
}

type point interface {
	X() float64
	Y() float64
}

func tileOffset(bb bbox) (xOffset, yOffset float64) {
	return bb.NW.X(), bb.NW.Y()
}

func tileScalingFactor(bb bbox, extent int) (xScale, yScale float64) {
	var (
		deltaX = math.Abs(bb.SE.X() - bb.NW.X())
		deltaY = math.Abs(bb.SE.Y() - bb.NW.Y())
	)
	return deltaX * float64(extent), deltaY * float64(extent)
}

func tileCoord(pt spatial.Point, extent int, xScale, yScale float64, xOffset, yOffset float64) (x, y int) {
	return int((pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent)),
		int((pt.Y() - yOffset) / (yScale / float64(extent)) * float64(extent))
}
