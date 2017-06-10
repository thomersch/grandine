package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineSegmentToCarthesian(t *testing.T) {
	a, b, c := LineSegmentToCarthesian(Point{0, 3}, Point{4, 1})
	assert.Equal(t, float64(2), a)
	assert.Equal(t, float64(4), b)
	assert.Equal(t, float64(12), c)
}
