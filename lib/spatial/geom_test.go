package spatial

import (
	"bytes"
	"encoding/json"
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
	_, err = wkb.Unmarshal(buf)
	assert.Nil(t, err)

	// test against own implementation
	rp := &Geom{}
	err = rp.UnmarshalWKB(bytes.NewReader(buf))
	assert.Nil(t, err)
	pt, err := rp.Point()
	assert.Nil(t, err)
	assert.Equal(t, spt, pt)
}

func TestMarshalWKBLineString(t *testing.T) {
	sls := []Point{{1, 2}, {3, 4}, {5, 4}}
	g, err := NewGeom(sls)
	assert.Nil(t, err)
	buf, err := g.MarshalWKB()
	assert.Nil(t, err)

	_, err = wkb.Unmarshal(buf)
	assert.Nil(t, err)

	rp := &Geom{}
	err = rp.UnmarshalWKB(bytes.NewReader(buf))
	assert.Nil(t, err)
	ls, err := rp.LineString()
	assert.Nil(t, err)
	assert.Equal(t, sls, ls)
}

func TestMarshalWKBPolygon(t *testing.T) {
	spoly := [][]Point{
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
	b.ReportAllocs()
	g, err := NewGeom(Point{2, 3})
	assert.Nil(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.MarshalWKB()
	}
}

func BenchmarkWKBMarshalRawPoint(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	p := Point{2, 3}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wkbWritePoint(&buf, p)
	}
}

func BenchmarkWKBMarshalLineString(b *testing.B) {
	b.ReportAllocs()
	ls := []Point{{2, 3}, {5, 6}, {10, 15}, {20, 50}}
	g, err := NewGeom(ls)
	assert.Nil(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.MarshalWKB()
	}
}

func BenchmarkWKBMarshalPoly(b *testing.B) {
	b.ReportAllocs()
	poly := [][]Point{{{2, 3}, {5, 6}, {10, 15}, {2, 3}}, {{10, 15}, {5, 6}, {10, 15}}}
	g, err := NewGeom(poly)
	assert.Nil(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.MarshalWKB()
	}
}

func BenchmarkWKBMarshalRawPoly(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	poly := [][]Point{{{2, 3}, {5, 6}, {10, 15}, {2, 3}}, {{10, 15}, {5, 6}, {10, 15}}}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wkbWritePolygon(&buf, poly)
	}
}
