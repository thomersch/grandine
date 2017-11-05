package mvt

import (
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/stretchr/testify/assert"
)

func TestScalePoint(t *testing.T) {
	bb := spatial.BBox{
		spatial.Point{50, 10},
		spatial.Point{52, 12},
	}
	ext := 4096
	xScale, yScale := tileScalingFactor(bb, ext)
	xOffset, yOffset := tileOffset(bb)

	tX, tY := tileCoord(spatial.Point{50, 10}, ext, xScale, yScale, xOffset, yOffset)
	assert.Equal(t, 0, tX)
	assert.Equal(t, ext, tY)

	tX, tY = tileCoord(spatial.Point{51, 10}, ext, xScale, yScale, xOffset, yOffset)
	assert.Equal(t, ext/2, tX)
	assert.Equal(t, ext, tY)

	tX, tY = tileCoord(spatial.Point{52, 12}, ext, xScale, yScale, xOffset, yOffset)
	assert.Equal(t, ext, tX)
	assert.Equal(t, 0, tY)
}

func TestProj4326To3857(t *testing.T) {
	assert.Equal(t, spatial.Point{4.57523107160354e+06, 2.28488107006733e+06}, proj4326To3857(spatial.Point{41.1, 20.1}).RoundedCoords())
	assert.Equal(t, spatial.Point{4.57523107160354e+06, -2.28488107006733e+06}, proj4326To3857(spatial.Point{41.1, -20.1}).RoundedCoords())
}

func scalePointToTileBarePoint(pt spatial.Point, extent int, xScale, yScale float64, xOffset, yOffset float64) spatial.Point {
	return spatial.Point{
		(pt.X - xOffset) / (xScale / float64(extent)) * float64(extent),
		(pt.Y - yOffset) / (yScale / float64(extent)) * float64(extent),
	}
}

func BenchmarkPointBare(b *testing.B) {
	pt := spatial.Point{1, 2}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		scalePointToTileBarePoint(pt, 4096, 1000, 1000, 3, 6)
	}
}
