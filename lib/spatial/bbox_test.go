package spatial

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBBoxIn(t *testing.T) {
	for n, tc := range []struct {
		inner BBox
		outer BBox
		in    bool
	}{
		{
			inner: BBox{
				Point{20, 20}, Point{30, 30},
			},
			outer: BBox{
				Point{10, 10}, Point{50, 50},
			},
			in: true,
		},
		{
			inner: BBox{
				Point{10, 10}, Point{50, 50},
			},
			outer: BBox{
				Point{20, 20}, Point{30, 30},
			},
			in: true,
		},
		{
			inner: BBox{
				Point{30, 10}, Point{40, 20},
			},
			outer: BBox{
				Point{45, 10}, Point{95, 60},
			},
			in: false,
		},
		{
			inner: BBox{
				Point{70, 10}, Point{80, 20},
			},
			outer: BBox{
				Point{10, 10}, Point{60, 60},
			},
			in: false,
		},
		{
			inner: BBox{
				Point{70, 80}, Point{95, 95},
			},
			outer: BBox{
				Point{10, 10}, Point{60, 60},
			},
			in: false,
		},
	} {
		t.Run(fmt.Sprintf("%v", n), func(t *testing.T) {
			assert.Equal(t, tc.in, tc.inner.Overlaps(tc.outer))
		})
	}
}

func TestBBoxFullyIn(t *testing.T) {
	b1 := BBox{SW: Point{1, 1}, NE: Point{2, 2}}
	b2 := BBox{SW: Point{0, 0}, NE: Point{3, 3}}

	assert.True(t, b1.FullyIn(b2))
	assert.False(t, b2.FullyIn(b1))
}
