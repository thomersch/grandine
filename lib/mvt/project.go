package mvt

import (
	"log"
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

type bbox struct {
	SW, NE point
}

type point interface {
	X() float64
	Y() float64
}

func tileOffset(bb bbox) (xOffset, yOffset float64) {
	return bb.SW.X(), bb.SW.Y()
}

func tileScalingFactor(bb bbox, extent int) (xScale, yScale float64) {
	var (
		// TODO: check if delta is correct if tile crosses hemispheres
		deltaX = math.Abs(bb.SW.X() - bb.NE.X())
		deltaY = math.Abs(bb.SW.Y() - bb.NE.Y())
	)
	log.Printf("deltae: %v %v", deltaX, deltaY)
	return deltaX * float64(extent), deltaY * float64(extent)
}

func tileCoord(pt spatial.Point, extent int, xScale, yScale float64, xOffset, yOffset float64) (x, y int) {
	return int((pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent)),
		int((pt.Y() - yOffset) / (yScale / float64(extent)) * float64(extent))
}
