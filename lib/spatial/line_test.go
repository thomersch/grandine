package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineIntersect(t *testing.T) {
	l1 := Segment{
		{0, 0},
		{0, 2},
	}
	l2 := Segment{
		{0, 3},
		{0, 4},
	}
	l3 := Segment{
		{-1, 1},
		{1, 1},
	}
	l4 := Segment{
		{1, 1},
		{2, 1},
	}

	t.Run("perpendicular intersect", func(t *testing.T) {
		isp, is := l1.Intersection(l3)
		assert.True(t, is)
		assert.Equal(t, Point{0, 1}, isp)
	})
	t.Run("non-intersect", func(t *testing.T) {
		_, is := l1.Intersection(l2)
		assert.False(t, is)
	})
	t.Run("non-segment intersect", func(t *testing.T) {
		isp, is := l1.Intersection(l4)
		assert.Equal(t, Point{0, 1}, isp)
		assert.False(t, is)
	})
}

func TestBBoxBorders(t *testing.T) {
	brds := BBoxBorders(Point{0, 0}, Point{0.5, 1})
	assert.Equal(t, Segment{{0, 0}, {0, 1}}, brds[0])
	assert.Equal(t, Segment{{0, 1}, {0.5, 1}}, brds[1])
	assert.Equal(t, Segment{{0.5, 1}, {0.5, 0}}, brds[2])
	assert.Equal(t, Segment{{0.5, 0}, {0, 0}}, brds[3])
}

func TestLineBBoxIntersect(t *testing.T) {
	l1 := Segment{
		{0, 0},
		{1, 0},
	}
	bbox := BBoxBorders(Point{0, 0}, Point{0.5, 1})

	t.Run("half cut", func(t *testing.T) {
		isp, is := l1.Intersection(bbox[2])
		assert.Equal(t, Segment{{0.5, 1}, {0.5, 0}}, bbox[2])
		assert.True(t, is)
		assert.Equal(t, Point{0.5, 0}, isp)
	})
}

func TestLineSegments(t *testing.T) {
	ln := Line{
		{0, 1},
		{1, 1},
		{1, 0},
	}
	segs := ln.Segments()
	assert.Equal(t, []Segment{
		{{0, 1}, {1, 1}},
		{{1, 1}, {1, 0}},
	}, segs)

	assert.Equal(t, ln, NewLinesFromSegments(segs)[0])
}

func TestLineSegmentsGapped(t *testing.T) {
	segs := []Segment{
		{{0, 1}, {1, 1}},
		{{1, 1}, {1, 2}},
		{{2, 2}, {3, 3}},
		{{3, 3}, {4, 3}},
		{{4, 3}, {5, 6}},
	}
	assert.Equal(t, []Line{
		{{0, 1}, {1, 1}, {1, 2}},
		{{2, 2}, {3, 3}, {4, 3}, {5, 6}},
	}, NewLinesFromSegments(segs))
}

func TestSegmentSplitAt(t *testing.T) {
	seg := Segment{
		{0, 0},
		{0, 5},
	}
	ss1, ss2 := seg.SplitAt(Point{0, 2})
	assert.Equal(t, Segment{{0, 0}, {0, 2}}, ss1)
	assert.Equal(t, Segment{{0, 2}, {0, 5}}, ss2)
}

func TestSegmentFullyIn(t *testing.T) {
	seg := Segment{
		{0, 0},
		{0, 5},
	}
	assert.True(t, seg.FullyInBBox(Point{-5, -5}, Point{5, 5}))
	assert.False(t, seg.FullyInBBox(Point{1, 0}, Point{1, 5}))
}

func BenchmarkLineIntersect(b *testing.B) {
	l1 := Segment{
		{0, 0},
		{0, 2},
	}
	l2 := Segment{
		{-1, 1},
		{1, 1},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l1.Intersection(l2)
	}
}

func TestLineBBox(t *testing.T) {
	l := Line{
		{5, 4},
		{2, 9},
		{5, 4},
		{-25, 4},
	}

	ne, sw := l.BBox()
	assert.Equal(t, Point{-25, 4}, ne)
	assert.Equal(t, Point{5, 9}, sw)
}

func TestClipLineString(t *testing.T) {
	ls1, err := NewGeom(Line{
		{1, 1},
		{1, 2},
		{2, 2},
		{3, 3},
	})
	assert.Nil(t, err)
	t.Run("completely inside bbox", func(t *testing.T) {
		assert.Equal(t, []Geom{ls1}, ls1.ClipToBBox(Point{0, 0}, Point{3, 3}))
	})
	t.Run("completely outside 1", func(t *testing.T) {
		assert.Equal(t, []Geom{}, ls1.ClipToBBox(Point{5, 5}, Point{12, 10}))
	})
	t.Run("completely outside 2", func(t *testing.T) {
		assert.Equal(t, []Geom{}, ls1.ClipToBBox(Point{-5, -5}, Point{0, 0}))
	})

	ls2, err := NewGeom([]Point{
		{1, 1},
		{3, 3},
		{5, 1},
	})
	assert.Nil(t, err)
	t.Run("split into two sublines", func(t *testing.T) {
		sl1, err := NewGeom([]Point{
			{1, 1},
			{2, 2},
		})
		assert.Nil(t, err)
		sl2, err := NewGeom([]Point{
			{4, 2},
			{5, 1},
		})
		assert.Nil(t, err)
		assert.Equal(t, []Geom{sl1, sl2}, ls2.ClipToBBox(Point{1, 1}, Point{5, 2}))
	})

	ls3, err := NewGeom(Line{
		{1, 1},
		{1, 2},
		{1, 5},
	})
	assert.Nil(t, err)
	t.Run("cut linestring", func(t *testing.T) {
		assert.Equal(t,
			[]Geom{{typ: GeomTypeLineString, g: Line{{1, 1}, {1, 2}, {1, 3}}}},
			ls3.ClipToBBox(Point{0, 0}, Point{3, 3}))
	})

	ls4, err := NewGeom(Line{
		{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0},
	})
	t.Run("closed line cut", func(t *testing.T) {
		cut1, err := NewGeom(Line{
			{0, 0},
			{0.5, 0},
		})
		assert.Nil(t, err)
		cut2, err := NewGeom(Line{
			{0.5, 1},
			{0, 1},
			{0, 0},
		})

		assert.Equal(t, []Geom{cut1, cut2}, ls4.ClipToBBox(Point{0, 0}, Point{0.5, 1}))
	})
}
