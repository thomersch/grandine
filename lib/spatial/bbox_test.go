package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBBox(t *testing.T) {
	r := []Point{
		{5, 4},
		{2, 9},
		{5, 4},
		{-25, 4},
	}

	ne, sw := ringBBox(r)
	assert.Equal(t, Point{-25, 4}, ne)
	assert.Equal(t, Point{5, 9}, sw)
}
