package mvt

import (
	"math"

	"github.com/thomersch/grandine/lib/spatial"
)

const earthRadius = 6378137

// flip axis
func flip(v int, extent int) int {
	if v == 0 {
		return extent
	}
	return (extent - v) % extent
}

func tileOffset(bb spatial.BBox) (xOffset, yOffset float64) {
	p := proj4326To3857(bb.SW)
	return p.X, p.Y
}

func tileScalingFactor(bb spatial.BBox, extent int) (xScale, yScale float64) {
	var (
		pSW    = proj4326To3857(bb.SW)
		pNE    = proj4326To3857(bb.NE)
		deltaX = math.Abs(pSW.X - pNE.X)
		deltaY = math.Abs(pSW.Y - pNE.Y)
	)
	return deltaX * float64(extent), deltaY * float64(extent)
}

func proj4326To3857(pt spatial.Point) spatial.Point {
	return spatial.Point{
		/* Lon/X */ degToRad(pt.X) * earthRadius,
		/* Lat/Y */ math.Log(math.Tan(degToRad(pt.Y)/2+math.Pi/4)) * earthRadius,
	}
}

func tileCoord(p spatial.Point, tp tileParams) (x, y int) {
	pt := proj4326To3857(p)
	return int((pt.X - tp.xOffset) / (tp.xScale / float64(tp.extent)) * float64(tp.extent)),
		flip(int((pt.Y-tp.yOffset)/(tp.yScale/float64(tp.extent))*float64(tp.extent)), tp.extent)
}

func degToRad(v float64) float64 {
	return v / (180 / math.Pi)
}

func radToDeg(v float64) float64 {
	return v * (180 / math.Pi)
}

type tileParams struct {
	xScale, yScale   float64
	xOffset, yOffset float64
	extent           int
}
