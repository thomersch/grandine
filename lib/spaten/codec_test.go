package spaten

import (
	"bytes"
	"testing"

	"github.com/thomersch/grandine/lib/geojson"
	"github.com/thomersch/grandine/lib/spatial"
)

func BenchmarkCodecThroughput(b *testing.B) {
	var (
		fc  = &spatial.FeatureCollection{Features: []spatial.Feature{}}
		sc  = &Codec{}
		gjc = &geojson.Codec{}
		buf = bytes.NewBuffer(make([]byte, 0, 200000000))
	)
	for i := 0; i < 50000; i++ {
		fc.Features = append(fc.Features, spatial.Feature{
			Geometry: spatial.MustNewGeom(spatial.Point{1, 2}),
			Props:    map[string]interface{}{"weight": 0},
		})
	}
	for i := 0; i < 50000; i++ {
		fc.Features = append(fc.Features, spatial.Feature{
			Geometry: spatial.MustNewGeom([]spatial.Point{{1, 2}, {3, 5}, {9, 0}, {2, 9}}),
			Props:    map[string]interface{}{"value": 14, "description": "i am a line"},
		})
	}

	b.Run("Spaten", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			sc.Encode(buf, fc)
		}
	})
	b.Run("GeoJSON", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			gjc.Encode(buf, fc)
		}
	})
}
