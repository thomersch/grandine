package spatial

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkClipToBBox(b *testing.B) {
	f, err := os.Open("testfiles/polygon_with_holes.geojson")
	assert.Nil(b, err)

	var fc FeatureCollection
	err = json.NewDecoder(f).Decode(&fc)
	assert.Nil(b, err)
	assert.Equal(b, 1, len(fc.Features))

	poly, err := fc.Features[0].Geometry.Polygon()
	assert.Nil(b, err)
	bbox := BBox{SW: Point{27.377929, 60.930432}, NE: Point{29.53125, 62.754725}}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		poly.ClipToBBox(bbox)
	}
}

func BenchmarkStringRepr(b *testing.B) {
	p := Polygon{
		Line{
			Point{1, 2}, Point{3, 4},
		},
		Line{
			Point{1, 2}, Point{3, 4},
		},
		Line{
			Point{1, 2}, Point{3, 4},
		},
	} // this is probably not valid, but this is not important for that benchmark

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.String()
	}
}
