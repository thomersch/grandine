package spaten

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
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

	err := WriteBlock(&buf, fcoll.Features, nil)
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

	err := WriteBlock(&buf, fs, nil)
	assert.Nil(t, err)

	const headerLength = 8 // TODO: consider exporting this
	// Compare buffer size with size written in header.
	assert.Equal(t, buf.Len()-headerLength, int(binary.LittleEndian.Uint32(buf.Bytes()[:4])))
	assert.Equal(t, "00000000", fmt.Sprintf("%x", buf.Bytes()[4:8]))
}

func TestInvalidBlockSize(t *testing.T) {
	buf, err := hex.DecodeString("FFFFFFFF00000000AAAA")
	assert.Nil(t, err)

	fc := spatial.NewFeatureCollection()
	err = readBlock(bytes.NewBuffer(buf), fc)
	assert.NotNil(t, err)
}

func TestWeirdFiles(t *testing.T) {
	var fls = []struct {
		buf       string
		shouldErr bool
	}{
		{"53504154000000000000000000000a0012171a15010100000000000000002440e523e8ca28c5517c1df8aa9998c44a40", true},
		{"53504154000000000000000000000000", false},
	}

	for i, f := range fls {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var c Codec
			buf, err := hex.DecodeString(f.buf)
			assert.Nil(t, err)

			fc := spatial.NewFeatureCollection()
			err = c.Decode(bytes.NewBuffer(buf), fc)
			if f.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
