package tile

import (
	"testing"

	"github.com/thomersch/grandine/lib/spatial"
)

func TestCoverage(t *testing.T) {
	Coverage(spatial.BBox{spatial.Point{-5, -5}, spatial.Point{10, 10}}, 7)
}
