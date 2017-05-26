package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClipPolygon(t *testing.T) {
	poly1, err := NewGeom(Polygon{
		{
			{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0},
		},
	})
	assert.Nil(t, err)
	t.Run("uncut", func(t *testing.T) {
		assert.Equal(t, []Geom{poly1}, poly1.ClipToBBox(Point{0, 0}, Point{1, 1}))
	})
	// t.Run("cut with single ring", func(t *testing.T) {
	// 	polyCut, err := NewGeom(Polygon{
	// 		{
	// 			{0, 0},
	// 			{0.5, 0},
	// 			{0.5, 1},
	// 			{0, 1},
	// 			{0, 0},
	// 		},
	// 	})
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, []Geom{polyCut}, poly1.ClipToBBox(Point{0, 0}, Point{0.5, 1}))
	// })
}
