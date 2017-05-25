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
