package spatial

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeWKBNullLineString(t *testing.T) {
	b, _ := hex.DecodeString("010300000030303000000000003030303030303030")
	buf := bytes.NewBuffer(b)

	var g Geom
	err := g.UnmarshalWKB(buf)
	assert.NotNil(t, err)
}

func TestGeomFromWKB(t *testing.T) {
	f, err := os.Open("testfiles/polygon.wkb")
	assert.Nil(t, err)
	defer f.Close()

	g, err := GeomFromWKB(f)
	assert.Nil(t, err)
	assert.Equal(t, g.Typ(), GeomTypePolygon)
}

func BenchmarkUnmarshalWKB(b *testing.B) {
	buf, _ := hex.DecodeString("03000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001040")

	b.Run("old style", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var g Geom
			g.UnmarshalWKB(bytes.NewBuffer(buf))
		}
	})

	b.Run("new style", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			GeomFromWKB(bytes.NewBuffer(buf))
		}
	})
}
