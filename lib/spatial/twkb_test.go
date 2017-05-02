package spatial

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTWKBWritePoint(t *testing.T) {
	origPt := Point{-212, 12.3}
	buf := bytes.Buffer{}
	err := twkbWritePoint(&buf, origPt, Point{})
	assert.Nil(t, err)
	fmt.Printf("%x\n", buf.Bytes())

	pt, err := twkbReadPoint(&buf, Point{})
	assert.Nil(t, err)
	assert.Equal(t, origPt, pt)
}

func BenchmarkTWKBWriteRawPoint(b *testing.B) {
	p := Point{2, 3}
	buf := bytes.Buffer{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := twkbWritePoint(&buf, p, Point{})
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
		twkbReadPoint(r, Point{})
	}
}
