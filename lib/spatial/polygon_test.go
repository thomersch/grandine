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
	poly2, err := NewGeom(Polygon{
		{
			{0, 0}, {0, 0.2}, {0.8, 0.2}, {0.8, 0.8}, {0, 0.8}, {0, 1}, {1, 1}, {1, 0}, {0, 0},
		},
	})
	assert.Nil(t, err)
	t.Run("uncut", func(t *testing.T) {
		assert.Equal(t, []Geom{poly1}, poly1.ClipToBBox(Point{0, 0}, Point{1, 1}))
	})
	t.Run("single ring cut", func(t *testing.T) {
		polyCut, err := NewGeom(Polygon{
			{
				{0.5, 1},
				{0, 1},
				{0, 0},
				{0.5, 0},
				{0.5, 1},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, []Geom{polyCut}, poly1.ClipToBBox(Point{0, 0}, Point{0.5, 1}))
	})
	t.Run("single ring into two subpolygons", func(t *testing.T) {
		polyCut1, err := NewGeom(Polygon{
			{
				{0.5, 0},
				{0, 0},
				{0, 0.2},
				{0.5, 0.2},
				{0.5, 0},
			},
		})
		assert.Nil(t, err)
		polyCut2, err := NewGeom(Polygon{
			{
				{0.5, 0.8},
				{0, 0.8},
				{0, 1},
				{0.5, 1},
				{0.5, 0.8},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, []Geom{polyCut1, polyCut2}, poly2.ClipToBBox(Point{0, 0}, Point{0.5, 1}))
	})
}

func TestPolygonsFromLines(t *testing.T) {
	ln := []Line{
		{{0, 0}, {0.5, 0}},
		{{0.5, 1}, {0, 1}, {0, 0}},
	}

	assert.Equal(t, []Polygon{{
		{
			{0.5, 1}, {0, 1}, {0, 0}, {0.5, 0}, {0.5, 1},
		},
	}}, PolygonsFromLines(ln))
}
