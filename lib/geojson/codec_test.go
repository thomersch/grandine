package geojson

import (
	"bytes"
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
	assert.Equal(t, "4326", fc.SRID)
	assert.Len(t, fc.Features, 1)
}

func TestDecodeID(t *testing.T) {
	f, err := os.Open("testdata/id.geojson")
	assert.Nil(t, err)

	var (
		c  = &Codec{}
		fc = spatial.FeatureCollection{}
	)
	err = c.Decode(f, &fc)
	assert.Nil(t, err)
	assert.Len(t, fc.Features, 2)
	assert.Equal(t, fc.Features[0].Properties()["id"], "asdf")
	assert.NotContains(t, fc.Features[1].Properties(), "id")
}

func TestDecodeMultipolygon(t *testing.T) {
	f, err := os.Open("testdata/multipolygon.geojson")
	assert.Nil(t, err)

	var (
		c  = &Codec{}
		fc = spatial.FeatureCollection{}
	)
	err = c.Decode(f, &fc)
	assert.Nil(t, err)
	assert.Len(t, fc.Features, 2)
}

func TestEncode(t *testing.T) {
	var (
		fc  spatial.FeatureCollection
		c   Codec
		buf = bytes.NewBuffer(make([]byte, 1000))
	)

	fc.SRID = "4326"
	err := c.Encode(buf, &fc)
	assert.Nil(t, err)
}
