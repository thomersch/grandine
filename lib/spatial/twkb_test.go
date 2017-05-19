package spatial

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTWKBReadHeader(t *testing.T) {
	buf, err := hex.DecodeString("24FF")
	assert.Nil(t, err)
	r := bytes.NewBuffer(buf)
	hd, err := twkbReadHeader(r)
	assert.Nil(t, err)
	assert.True(t, hd.bbox)
}

func TestTWKBWriteHeader(t *testing.T) {
	w := &bytes.Buffer{}
	typ := GeomTypeLineString
	precision := 4
	twkbWriteHeader(w, typ, precision)

	hd, err := twkbReadHeader(w)
	assert.Nil(t, err)
	assert.Equal(t, twkbHeader{typ: typ, precision: precision}, hd)
}

func TestTWKBWritePoint(t *testing.T) {
	precision := 6
	origPt := Point{-212, 12.3}
	buf := bytes.Buffer{}
	err := twkbWritePoint(&buf, origPt, Point{}, precision)
	assert.Nil(t, err)

	pt, err := twkbReadPoint(&buf, Point{}, precision)
	assert.Nil(t, err)
	assert.Equal(t, origPt, pt)
}

func TestTWKBReadPoint(t *testing.T) {
	buf, err := hex.DecodeString("01000204")
	assert.Nil(t, err)
	r := bytes.NewBuffer(buf)

	hd, err := twkbReadHeader(r)
	assert.Nil(t, err)
	pt, err := twkbReadPoint(r, Point{}, hd.precision)
	assert.Nil(t, err)
	assert.Equal(t, Point{1, 2}, pt)
}

func TestTWKBReadLine(t *testing.T) {
	buf, err := hex.DecodeString("02000202020808")
	assert.Nil(t, err)
	r := bytes.NewBuffer(buf)

	hd, err := twkbReadHeader(r)
	assert.Nil(t, err)
	ls, err := twkbReadLineString(r, hd.precision)
	assert.Nil(t, err)
	assert.Equal(t, []Point{{1, 1}, {5, 5}}, ls)
}

func BenchmarkTWKBWriteRawPoint(b *testing.B) {
	precision := 2
	p := Point{2, 3}
	buf := bytes.Buffer{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := twkbWritePoint(&buf, p, Point{}, precision)
		assert.Nil(b, err)
	}
}

func BenchmarkTWKBReadRawPoint(b *testing.B) {
	var rawPt []byte
	_, err := fmt.Sscanf("fff396ca01c0bbdd0b", "%x", &rawPt)
	assert.Nil(b, err)
	r := bytes.NewReader(rawPt)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Reset(rawPt)
		twkbReadPoint(r, Point{}, 0)
	}
}
