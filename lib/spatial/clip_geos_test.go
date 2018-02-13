// +build !golangclip

package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeosConversion(t *testing.T) {
	ply := Polygon{{{0, 1}, {2.5, 1}, {3, 2}, {2.5, 3}, {1, 3}}}
	gs := ply.geos()
	assert.Equal(t, ply, geosToPolygon(gs))
}
