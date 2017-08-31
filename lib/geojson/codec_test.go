package geojson

import (
	"os"
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("testdata/01.geojson")
	assert.Nil(t, err)

	var (
		c  = &Codec{}
		fc = spatial.FeatureCollection{}
	)
	err = c.Decode(f, &fc)
	assert.Nil(t, err)
	assert.Equal(t, fc.SRID, "4326")
}
