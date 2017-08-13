package mvt

import (
	"log"
	"math"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/pebbe/go-proj-4/proj"
)

var proj3857, proj4326 *proj.Proj

func init() {
	var err error

	proj3857, err = proj.NewProj("+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs")
	if err != nil {
		panic("could not initialize proj")
	}
	proj4326, err = proj.NewProj("+proj=longlat +datum=WGS84 +no_defs")
	if err != nil {
		panic("could not initialize proj")
	}
}

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
	p := proj4326To3857(bb.SW)
	return p.X(), p.Y()
}

func tileScalingFactor(bb spatial.BBox, extent int) (xScale, yScale float64) {
	var (
		pSW    = proj4326To3857(bb.SW)
		pNE    = proj4326To3857(bb.NE)
		deltaX = math.Abs(pSW.X() - pNE.X())
		deltaY = math.Abs(pSW.Y() - pNE.Y())
	)
	return deltaX * float64(extent), deltaY * float64(extent)
}

func proj4326To3857(pt spatial.Point) spatial.Point {
	x, y, err := proj.Transform2(proj4326, proj3857, proj.DegToRad(pt.X()), proj.DegToRad(pt.Y()))
	if err != nil {
		log.Println(pt)
		panic(err)
	}
	return spatial.Point{x, y}
}

func tileCoord(p spatial.Point, extent int, xScale, yScale float64, xOffset, yOffset float64) (x, y int) {
	pt := proj4326To3857(p)
	return int((pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent)),
		flip(int((pt.Y()-yOffset)/(yScale/float64(extent))*float64(extent)), extent)
}
