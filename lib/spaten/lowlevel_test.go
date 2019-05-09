package spaten

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

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
		{"53504154000000001b00000030303030303012171a15010300000030303000000000003030303030303030", true},
		{"53504154000000001010101000000000", true},
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

func BenchmarkReadBlock(b *testing.B) {
	var (
		buf = bytes.NewBuffer([]byte{})
		fs  = spatial.NewFeatureCollection()
	)
	err := WriteBlock(buf, fs.Features, nil)
	assert.Nil(b, err)
	r := bytes.NewReader(buf.Bytes())

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Seek(0, 0)
		err := readBlock(r, fs)
		assert.Nil(b, err)
	}
}

func BenchmarkReadBlockThroughput(b *testing.B) {
	var (
		wBuf = bytes.NewBuffer([]byte{})
		fs   = spatial.NewFeatureCollection()
	)

	err := WriteBlock(wBuf, []spatial.Feature{
		{
			Geometry: spatial.MustNewGeom(spatial.Point{2, 3}),
			Props: map[string]interface{}{
				"highway": "primary",
				"number":  1,
			},
		},
	}, nil)
	assert.Nil(b, err)

	var (
		ptBuf   = wBuf.Bytes()[8:]
		fullBuf = make([]byte, 8)
	)
	for n := 0; n < 100000; n++ {
		fullBuf = append(fullBuf, ptBuf...)
	}
	binary.LittleEndian.PutUint32(fullBuf[:4], uint32(len(fullBuf)-8))

	b.ReportAllocs()
	b.ResetTimer()

	t := time.Now()
	for n := 0; n < b.N; n++ {
		r := bytes.NewBuffer(fullBuf)
		err = readBlock(r, fs)
		assert.Nil(b, err)
		fs.Features = []spatial.Feature{}
	}
	b.Logf("%v bytes read, in %v blocks,		throughput: %v B/s", len(fullBuf)*b.N, b.N, int(float64(len(fullBuf)*b.N)/time.Since(t).Seconds()))
}
