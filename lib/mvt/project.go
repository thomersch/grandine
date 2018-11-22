package mvt

import (
	"math"

	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

const earthRadius = 6378137

// flip axis
func flip(v int, extent int) int {
	return extent - v
}

func flipFloat(v float64, extent int) float64 {
	return float64(flip(int(v), extent))
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
	newPt := tilePoint(p, tp)
	return int(newPt.X), int(newPt.Y)
}

func tilePoint(p spatial.Point, tp tileParams) spatial.Point {
	pt := proj4326To3857(p)
	return spatial.Point{
		X: (pt.X - tp.xOffset) / (tp.xScale / float64(tp.extent)) * float64(tp.extent),
		Y: flipFloat((pt.Y-tp.yOffset)/(tp.yScale/float64(tp.extent))*float64(tp.extent), tp.extent),
	}
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

func newTileParams(tid tile.ID, ext int) tileParams {
	var tp tileParams
	tp.xScale, tp.yScale = tileScalingFactor(tid.BBox(), ext)
	tp.xOffset, tp.yOffset = tileOffset(tid.BBox())
	tp.extent = ext
	return tp
}
