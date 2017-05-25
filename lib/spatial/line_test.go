package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineIntersect(t *testing.T) {
	l1 := [2]Point{
		{0, 0},
		{0, 2},
	}
	l2 := [2]Point{
		{0, 3},
		{0, 4},
	}
	l3 := [2]Point{
		{-1, 1},
		{1, 1},
	}
	l4 := [2]Point{
		{1, 1},
		{2, 1},
	}

	t.Run("perpendicular intersect", func(t *testing.T) {
		isp, is := LineIntersect(l1, l3)
		assert.True(t, is)
		assert.Equal(t, Point{0, 1}, isp)
	})
	t.Run("non-intersect", func(t *testing.T) {
		_, is := LineIntersect(l1, l2)
		assert.False(t, is)
	})
	t.Run("non-segment intersect", func(t *testing.T) {
		isp, is := LineIntersect(l1, l4)
		assert.Equal(t, Point{0, 1}, isp)
		assert.False(t, is)
	})
}

func BenchmarkLineIntersect(b *testing.B) {
	l1 := [2]Point{
		{0, 0},
		{0, 2},
	}
	l2 := [2]Point{
		{-1, 1},
		{1, 1},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		LineIntersect(l1, l2)
	}
}
