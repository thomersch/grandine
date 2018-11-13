package spatial

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom/encoding/wkb"
)

func TestMarshalWKBPoint(t *testing.T) {
	spt := Point{-21.123456, 0.981231}
	g, err := NewGeom(spt)
	assert.Nil(t, err)
	buf, err := g.MarshalWKB()
	assert.Nil(t, err)

	// test against third party implementation
	tpt, err := wkb.Unmarshal(buf)
	assert.Nil(t, err)
	assert.Equal(t, tpt.FlatCoords()[0], spt.X)
	assert.Equal(t, tpt.FlatCoords()[1], spt.Y)

	// test against own implementation
	rp := &Geom{}
	err = rp.UnmarshalWKB(bytes.NewReader(buf))
	assert.Nil(t, err)
	pt, err := rp.Point()
	assert.Nil(t, err)
	assert.Equal(t, spt, pt)
}

func TestWKBUnmarshal(t *testing.T) {
	f, err := os.Open("testfiles/polygon.wkb")
	assert.Nil(t, err)

	buf, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	var g Geom
	err = g.UnmarshalWKB(bytes.NewBuffer(buf))
	assert.Nil(t, err)
	assert.Equal(t, g.Typ(), GeomTypePolygon)

	bufOut, err := g.MarshalWKB()
	assert.Nil(t, err)
	assert.Equal(t, buf, bufOut)

	// TODO: fix the test, which broke because of winding order reconstruction
	// gj, err := os.Open("testfiles/polygon.json")
	// assert.Nil(t, err)
	// gjbuf, err := ioutil.ReadAll(gj)

	// var g2 Geom
	// g2.UnmarshalJSON(gjbuf)
	// assert.Equal(t, g, g2)
}

func TestUnmarshalWKBEOF(t *testing.T) {
	var buf []byte
	fmt.Sscanf("09000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001040", "%x", &buf)

	_, err := wkbReadLineString(bytes.NewReader(buf))
	assert.Equal(t, io.EOF, err)
}

func TestMarshalWKBLineString(t *testing.T) {
	sls := Line{{1, 2}, {3, 4}, {5, 4}}
	g, err := NewGeom(sls)
	assert.Nil(t, err)
	buf, err := g.MarshalWKB()
	assert.Nil(t, err)

	tls, err := wkb.Unmarshal(buf)
	assert.Nil(t, err)
	assert.Equal(t, tls.FlatCoords()[0], sls[0].X)
	assert.Equal(t, tls.FlatCoords()[1], sls[0].Y)
	assert.Equal(t, tls.FlatCoords()[2], sls[1].X)
	assert.Equal(t, tls.FlatCoords()[3], sls[1].Y)
	assert.Equal(t, tls.FlatCoords()[4], sls[2].X)
	assert.Equal(t, tls.FlatCoords()[5], sls[2].Y)

	rp := &Geom{}
	err = rp.UnmarshalWKB(bytes.NewReader(buf))
	assert.Nil(t, err)
	ls, err := rp.LineString()
	assert.Nil(t, err)
	assert.Equal(t, sls, ls)
}

func TestMarshalWKBPolygon(t *testing.T) {
	spoly := Polygon{
		{
			{1, 2}, {3, 4}, {5, 4},
		},
		{
			{2, 2}, {3, 4}, {2, 2},
		},
	}
	g, err := NewGeom(spoly)
	assert.Nil(t, err)
	buf, err := g.MarshalWKB()
	assert.Nil(t, err)

	_, err = wkb.Unmarshal(buf)
	assert.Nil(t, err)

	rp := &Geom{}
	err = rp.UnmarshalWKB(bytes.NewReader(buf))
	assert.Nil(t, err)
	poly, err := rp.Polygon()
	assert.Nil(t, err)
	assert.Equal(t, spoly, poly)
}

func TestGeoJSON(t *testing.T) {
	f, err := os.Open("testfiles/featurecollection.geojson")
	assert.Nil(t, err)
	defer f.Close()

	fc := FeatureCollection{}
	err = json.NewDecoder(f).Decode(&fc)
	assert.Nil(t, err)

	p, err := fc.Features[0].Geometry.Point()
	assert.Nil(t, err)
	assert.NotNil(t, p)

	ls, err := fc.Features[1].Geometry.LineString()
	assert.Nil(t, err)
	assert.NotNil(t, ls)

	poly, err := fc.Features[2].Geometry.Polygon()
	assert.Nil(t, err)
	assert.NotNil(t, poly)

	buf, err := json.Marshal(fc)
	assert.Nil(t, err)
	assert.NotNil(t, buf)
}

func BenchmarkWKBMarshalPoint(b *testing.B) {
	g, err := NewGeom(Point{2, 3})
	assert.Nil(b, err)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.MarshalWKB()
	}
}

func BenchmarkWKBMarshalRawPoint(b *testing.B) {
	var buf bytes.Buffer
	p := Point{2, 3}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wkbWritePoint(&buf, p)
	}
}

func BenchmarkWKBMarshalLineString(b *testing.B) {
	ls := []Point{{2, 3}, {5, 6}, {10, 15}, {20, 50}}
	g, err := NewGeom(ls)
	assert.Nil(b, err)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.MarshalWKB()
	}
}

func BenchmarkWKBMarshalPoly(b *testing.B) {
	poly := [][]Point{{{2, 3}, {5, 6}, {10, 15}, {2, 3}}, {{10, 15}, {5, 6}, {10, 15}}}
	g, err := NewGeom(poly)
	assert.Nil(b, err)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.MarshalWKB()
	}
}

func BenchmarkWKBMarshalRawPoly(b *testing.B) {
	var buf bytes.Buffer
	poly := Polygon{{{2, 3}, {5, 6}, {10, 15}, {2, 3}}, {{10, 15}, {5, 6}, {10, 15}}}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wkbWritePolygon(&buf, poly)
	}
}

func BenchmarkWKBUnmarshalPoint(b *testing.B) {
	var rawPt []byte
	_, err := fmt.Sscanf("b77efacf9a1f35c0b648da8d3e66ef3f", "%x", &rawPt)
	assert.Nil(b, err)
	r := bytes.NewReader(rawPt)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Reset(rawPt)
		wkbReadPoint(r)
	}
}

func BenchmarkWKBUnmarshalLineString(b *testing.B) {
	var rawLine []byte
	_, err := fmt.Sscanf("03000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001040", "%x", &rawLine)
	assert.Nil(b, err)
	r := bytes.NewReader(rawLine)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Reset(rawLine)
		wkbReadLineString(r)
	}
}

func BenchmarkWKBUnmarshalPoly(b *testing.B) {
	var rawPoly []byte
	_, err := fmt.Sscanf("0200000003000000000000000000f03f0000000000000040000000000000084000000000000010400000000000001440000000000000104003000000000000000000004000000000000000400000000000000840000000000000104000000000000000400000000000000040", "%x", &rawPoly)
	assert.Nil(b, err)
	r := bytes.NewReader(rawPoly)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Reset(rawPoly)
		wkbReadPolygon(r)
	}
}

func TestClipPoint(t *testing.T) {
	p, err := NewGeom(Point{1, 1})
	assert.Nil(t, err)
	t.Run("inside", func(t *testing.T) {
		assert.Equal(t, []Geom{p}, p.ClipToBBox(BBox{Point{0, 0}, Point{5, 5}}))
	})
	t.Run("outside", func(t *testing.T) {
		assert.Equal(t, []Geom{}, p.ClipToBBox(BBox{Point{5, 0}, Point{5, 5}}))
	})
	t.Run("on SW edge", func(t *testing.T) {
		assert.Equal(t, []Geom{p}, p.ClipToBBox(BBox{Point{0, 0}, Point{1, 1}}))
	})
	t.Run("on NE edge", func(t *testing.T) {
		assert.Equal(t, []Geom{p}, p.ClipToBBox(BBox{Point{1, 1}, Point{2, 2}}))
	})
}
