package geojsonseq

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomersch/grandine/lib/spatial"
)

func TestChunkedDecode(t *testing.T) {
	f, err := os.Open("testdata/example.geojsonseq")
	assert.Nil(t, err)
	defer f.Close()

	c := Codec{}
	chunks, err := c.ChunkedDecode(f)
	assert.Nil(t, err)

	var fcoll spatial.FeatureCollection
	for chunks.Next() {
		err := chunks.Scan(&fcoll)
		assert.Nil(t, err)
	}
	assert.Len(t, fcoll.Features, 10)
}
