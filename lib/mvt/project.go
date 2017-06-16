package mvt

import (
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

type point interface {
	X() float64
	Y() float64
}

// flip axis
func flip(v int, extent int) int {
	if v == 0 {
		return extent
	}
	return (extent - v) % extent
}

func tileOffset(bb spatial.BBox) (xOffset, yOffset float64) {
	return bb.SW.X(), bb.SW.Y()
}

func tileScalingFactor(bb spatial.BBox, extent int) (xScale, yScale float64) {
	var (
		deltaX = math.Abs(bb.SW.X() - bb.NE.X())
		deltaY = math.Abs(bb.SW.Y() - bb.NE.Y())
	)
	return deltaX * float64(extent), deltaY * float64(extent)
}

func tileCoord(pt spatial.Point, extent int, xScale, yScale float64, xOffset, yOffset float64) (x, y int) {
	return int((pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent)),
		flip(int((pt.Y()-yOffset)/(yScale/float64(extent))*float64(extent)), extent)
}
