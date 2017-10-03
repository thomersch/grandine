// +build golangclip

package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClipPolygon(t *testing.T) {
	poly1, err := NewGeom(Polygon{
		{
			{0, 1}, {0, 0}, {1, 0}, {1, 1},
		},
	})
	poly2, err := NewGeom(Polygon{
		{
			{0, 0}, {0, 0.2}, {0.8, 0.2}, {0.8, 0.8}, {0, 0.8}, {0, 1}, {1, 1}, {1, 0},
		},
	})
	poly3, err := NewGeom(Polygon{
		{
			{0, 10}, {0, 0}, {10, 0},
		},
	})
	assert.Nil(t, err)
	t.Run("uncut", func(t *testing.T) {
		assert.Equal(t, []Geom{poly1}, poly1.ClipToBBox(BBox{Point{0, 0}, Point{1, 1}}))
	})
	t.Run("single ring cut", func(t *testing.T) {
		polyCut, err := NewGeom(Polygon{
			{
				{0, 1},
				{0, 0},
				{0.5, 0},
				{0.5, 1},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, []Geom{polyCut}, poly1.ClipToBBox(BBox{Point{0, 0}, Point{0.5, 1}}))
	})
	t.Run("single ring into two subpolygons", func(t *testing.T) {
		polyCut1, err := NewGeom(Polygon{
			{
				{0, 0.2},
				{0, 0},
				{0.5, 0},
				{0.5, 0.2},
			},
			{
				{0, 1},
				{0, 0.8},
				{0.5, 0.8},
				{0.5, 1},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, polyCut1, poly2.ClipToBBox(BBox{Point{-0.1, -0.1}, Point{0.5, 1.1}})[0])
	})
	t.Run("triangle cut", func(t *testing.T) {
		poly3.ClipToBBox(BBox{Point{5, -5}, Point{20, 20}})
	})

	// TODO: test cut where the bbox of the polygon overlaps with cut bbox, but isn't actually inside
}
