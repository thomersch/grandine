package mvt

import (
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/stretchr/testify/assert"
)

func TestScalePoint(t *testing.T) {
	bb := bbox{
		spatial.Point{50, 10},
		spatial.Point{52, 12},
	}
	ext := 4096
	xScale, yScale := tileScalingFactor(bb, ext)
	xOffset, yOffset := tileOffset(bb)

	tX, tY := tileCoord(spatial.Point{50, 10}, ext, xScale, yScale, xOffset, yOffset)
	assert.Equal(t, 0, tX)
	assert.Equal(t, 0, tY)

	tX, tY = tileCoord(spatial.Point{51, 10}, ext, xScale, yScale, xOffset, yOffset)
	assert.Equal(t, 2048, tX)
	assert.Equal(t, 0, tY)

	tX, tY = tileCoord(spatial.Point{52, 12}, ext, xScale, yScale, xOffset, yOffset)
	assert.Equal(t, 4096, tX)
	assert.Equal(t, 4096, tY)
}

func scalePointToTileBothInterface(pt point, extent int, xScale, yScale float64, xOffset, yOffset float64) point {
	return spatial.Point{
		(pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent),
		(pt.Y() - yOffset) / (yScale / float64(extent)) * float64(extent),
	}
}

func scalePointToTileBarePoint(pt spatial.Point, extent int, xScale, yScale float64, xOffset, yOffset float64) spatial.Point {
	return spatial.Point{
		(pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent),
		(pt.Y() - yOffset) / (yScale / float64(extent)) * float64(extent),
	}
}

func scalePointToTileInterfaceInput(pt point, extent int, xScale, yScale float64, xOffset, yOffset float64) spatial.Point {
	return spatial.Point{
		(pt.X() - xOffset) / (xScale / float64(extent)) * float64(extent),
		(pt.Y() - yOffset) / (yScale / float64(extent)) * float64(extent),
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

func BenchmarkPointInterfaceInputAndOutput(b *testing.B) {
	pt := spatial.Point{1, 2}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		scalePointToTileBothInterface(pt, 4096, 1000, 1000, 3, 6)
	}
}

func BenchmarkPointInterfaceInput(b *testing.B) {
	pt := spatial.Point{1, 2}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		scalePointToTileInterfaceInput(pt, 4096, 1000, 1000, 3, 6)
	}
}
