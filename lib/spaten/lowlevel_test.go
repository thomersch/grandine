package spaten

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomersch/grandine/lib/spatial"
)

// TestReadFileHeader is here in order to prevent regressions in header parsing.
func TestReadFileHeader(t *testing.T) {
	buf, err := hex.DecodeString("5350415400000000")
	assert.Nil(t, err)
	r := bytes.NewBuffer(buf)

	hd, err := ReadFileHeader(r)
	assert.Nil(t, err)
	assert.Equal(t, Header{Version: 0}, hd)
}

func TestHeaderSelfTest(t *testing.T) {
	var buf bytes.Buffer
	err := WriteFileHeader(&buf)
	assert.Nil(t, err)

	_, err = ReadFileHeader(&buf)
	assert.Nil(t, err)
}

func TestBlockSelfTest(t *testing.T) {
	var (
		buf   bytes.Buffer
		fcoll = spatial.FeatureCollection{
			Features: []spatial.Feature{
				{
					Props: map[string]interface{}{
						"key1": 1,
						"key2": "string",
						"key3": -12.981,
					},
					Geometry: spatial.MustNewGeom(spatial.Point{24, 1}),
				},
				{
					Props: map[string]interface{}{
						"yes": "NO",
					},
					Geometry: spatial.MustNewGeom(spatial.Line{{24, 1}, {25, 0}, {9, -4}}),
				},
				{
					Props: map[string]interface{}{
						"name": "RichardF Box",
					},
					Geometry: spatial.MustNewGeom(spatial.Polygon{{{24, 1}, {25, 0}, {9, -4}}}),
				},
			},
		}
	)

	err := WriteBlock(&buf, fcoll.Features)
	assert.Nil(t, err)

	var fcollRead spatial.FeatureCollection
	err = ReadBlocks(&buf, &fcollRead)
	assert.Nil(t, err)
	assert.Equal(t, fcoll, fcollRead)
}

func TestBlockHeaderEncoding(t *testing.T) {
	var (
		buf bytes.Buffer
		fs  = []spatial.Feature{
			{
				Geometry: spatial.MustNewGeom(spatial.Point{1, 2}),
			},
		}
	)

	err := WriteBlock(&buf, fs)
	// This is not the most robust test, but will fail if you accidentaly break the encoder.
	assert.Nil(t, err)
	assert.Equal(t, "1900000000000000", fmt.Sprintf("%x", buf.Bytes()[:8]))
	fmt.Printf("%x", buf.Bytes())
}
