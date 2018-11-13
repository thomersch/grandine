package spatial

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundWithPrecision(t *testing.T) {
	pt1 := Point{-5.4213000001, 10.9874000001}
	assert.Equal(t, Point{-5.4213, 10.9874}, pt1.RoundedCoords())
}

func TestPointInPolygon(t *testing.T) {
	square := Polygon{
		{
			{-1, 1}, {-1, -1}, {1, -1}, {1, 1},
		},
	}
	triangle := Polygon{
		{
			{0, 0}, {1, 2}, {2, 0},
		},
	}
	squareWithHole := Polygon{
		{
			{0, 0}, {0, 10}, {10, 10}, {10, 0},
		},
		{
			{2.5, 2.5}, {2.5, 7.5}, {7.5, 7.5}, {7.5, 2.5},
		},
	}

	t.Run("simple in", func(t *testing.T) {
		pt := Point{0, 0}
		assert.True(t, pt.InPolygon(square))
	})

	t.Run("simple out 1", func(t *testing.T) {
		pt := Point{-2, -2}
		assert.False(t, pt.InPolygon(square))
	})

	t.Run("simple out 2", func(t *testing.T) {
		pt := Point{3, 3}
		assert.False(t, pt.InPolygon(square))
	})

	t.Run("triangle in", func(t *testing.T) {
		pt := Point{1, 1}
		assert.True(t, pt.InPolygon(triangle))
	})
	t.Run("triangle out", func(t *testing.T) {
		pt := Point{0.5, 1.1}
		assert.False(t, pt.InPolygon(triangle))
	})

	t.Run("holed in", func(t *testing.T) {
		pt := Point{1, 1}
		assert.True(t, pt.InPolygon(squareWithHole))
	})
	t.Run("holed out", func(t *testing.T) {
		pt := Point{5, 5}
		assert.False(t, pt.InPolygon(squareWithHole))
	})

	t.Run("closing segment", func(t *testing.T) {
		pt := Point{25.48828125, -18.312810846425432}
		poly := Polygon{Line{Point{X: 7.3828125, Y: -23.241346102386135}, Point{X: 28.4765625, Y: -8.05922962720018}, Point{X: 55.1953125, Y: -11.178401873711772}, Point{X: 22.148437499999996, Y: -33.137551192346145}}}
		assert.True(t, pt.InPolygon(poly))
	})
}

func TestPointJSONMarshal(t *testing.T) {
	p := Point{-12.00000000001, 179.1}
	buf, err := json.Marshal(p)
	assert.Nil(t, err)
	fmt.Printf("%s", buf)
}

func BenchmarkPointJSONMarshal(b *testing.B) {
	p := Point{12.00000000001, 13.000000000001}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		json.Marshal(p)
	}
}
