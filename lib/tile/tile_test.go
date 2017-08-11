package tile

import (
	"fmt"
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/stretchr/testify/assert"
)

func TestTileName(t *testing.T) {
	for _, tc := range []struct {
		p        spatial.Point
		zl       int
		expected ID
	}{
		{
			p:        spatial.Point{13.73630, 51.05377},
			zl:       14,
			expected: ID{X: 8817, Y: 5481, Z: 14},
		},
		{
			p:        spatial.Point{18.39856, -33.90184},
			zl:       14,
			expected: ID{X: 9029, Y: 9833, Z: 14},
		},
		{
			p:        spatial.Point{-54.59123, -25.59547},
			zl:       14,
			expected: ID{X: 5707, Y: 9397, Z: 14},
		},
		{
			p:        spatial.Point{-21.94073, 64.14607},
			zl:       14,
			expected: ID{X: 7193, Y: 4354, Z: 14},
		},
		{
			p:        spatial.Point{-31.16580, 83.65691},
			zl:       14,
			expected: ID{X: 6773, Y: 648, Z: 14},
		},
		{
			p:        spatial.Point{-64.45649, -85.04438},
			zl:       14,
			expected: ID{X: 5258, Y: 16380, Z: 14},
		},
		{
			p:        spatial.Point{180, -90},
			zl:       1,
			expected: ID{X: 1, Y: 1, Z: 1},
		},
		{
			p:        spatial.Point{-180, 90},
			zl:       1,
			expected: ID{X: 0, Y: 0, Z: 1},
		},
	} {
		t.Run(fmt.Sprintf("%v_%v", tc.expected.X, tc.expected.Y), func(t *testing.T) {
			var fail bool
			ti := TileName(tc.p, tc.zl)
			if ti.X != tc.expected.X {
				fail = true
			}
			if ti.Y != tc.expected.Y {
				fail = true
			}
			if ti.Z != tc.expected.Z {
				fail = true
			}
			if fail {
				t.Fatalf("invalid conversion, expected %v, got %v", tc.expected, ti)
			}
		})
	}
}

func TestTileBBox(t *testing.T) {
	for _, tc := range []struct {
		tid      ID
		expected spatial.BBox
	}{
		{
			tid:      ID{0, 0, 0},
			expected: spatial.BBox{spatial.Point{-180, -85.05112878}, spatial.Point{180, 85.05112878}},
		},
		{
			tid:      ID{0, 0, 1},
			expected: spatial.BBox{spatial.Point{-180, 0}, spatial.Point{0, 85.05112878}},
		},
		{
			tid:      ID{0, 1, 1},
			expected: spatial.BBox{spatial.Point{-180, -85.05112878}, spatial.Point{0, 0}},
		},
		{
			tid:      ID{1, 2, 2},
			expected: spatial.BBox{spatial.Point{-90, -66.51326044}, spatial.Point{0, 0}},
		},
	} {
		t.Run(fmt.Sprintf("%v_%v_%v", tc.tid.X, tc.tid.Y, tc.tid.Z), func(t *testing.T) {
			bbox := tc.tid.BBox()
			bbox.NE = bbox.NE.RoundedCoords()
			bbox.SW = bbox.SW.RoundedCoords()
			assert.Equal(t, tc.expected, bbox)
		})
	}
}
