package spatial

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewind(t *testing.T) {
	p := Polygon{
		{{1, 2}, {8, 9}, {10, 12}, {1, 2}},
		{{0, 1}, {7, 9}, {0, 1}, {2, 12}, {0, 1}},
	}
	p.Rewind()
	assert.Equal(t, Polygon{
		{{1, 2}, {10, 12}, {8, 9}, {1, 2}},
		{{0, 1}, {2, 12}, {0, 1}, {7, 9}, {0, 1}},
	}, p)
}

func TestWinding(t *testing.T) {
	f, err := os.Open("testfiles/winding_wild.geojson")
	assert.Nil(t, err)
	defer f.Close()

	var fc FeatureCollection
	err = json.NewDecoder(f).Decode(&fc)
	assert.Nil(t, err)

	var outOrder []bool
	for _, ring := range fc.Features[0].Geometry.MustPolygon() {
		outOrder = append(outOrder, ring.Clockwise())
	}
	assert.Equal(t, []bool{true, false, true, true, false}, outOrder) // correct order
}

func TestFixWinding(t *testing.T) {
	g := Geom{typ: GeomTypePolygon, g: Polygon{
		Line{Point{X: -2.109375, Y: 11.178401873711785}, Point{X: -16.875, Y: -43.06888777416961}, Point{X: 62.57812500000001, Y: -43.580390855607845}, Point{X: 81.5625, Y: 8.407168163601076}},
		Line{Point{X: 7.3828125, Y: -23.241346102386135}, Point{X: 28.4765625, Y: -8.05922962720018}, Point{X: 55.1953125, Y: -11.178401873711772}, Point{X: 22.148437499999996, Y: -33.137551192346145}},
		Line{Point{X: 25.48828125, Y: -18.312810846425432}, Point{X: 33.22265625, Y: -16.720385051693988}, Point{X: 34.013671875, Y: -21.207458730482642}, Point{X: 23.466796875, Y: -24.766784522874428}},
		Line{Point{X: 27.5537109375, Y: -12.618897304044012}, Point{X: 29.02587890625, Y: -12.146745814539685}, Point{X: 29.377441406249996, Y: -14.604847155053898}, Point{X: 26.3671875, Y: -15.855673509998681}},
		Line{Point{X: 27.0703125, Y: -20.3034175184893}, Point{X: 27.509765625, Y: -21.616579336740593}, Point{X: 31.113281249999996, Y: -19.559790136497398}}}}

	poly := g.MustPolygon()
	var inOrder []bool
	for _, ring := range poly {
		inOrder = append(inOrder, ring.Clockwise())
	}
	assert.Equal(t, []bool{true, false, false, false, true}, inOrder) // wild order

	poly.FixWinding()

	var outOrder []bool
	for _, ring := range poly {
		outOrder = append(outOrder, ring.Clockwise())
	}
	assert.Equal(t, []bool{true, false, true, true, false}, outOrder) // correct order
}

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

func BenchmarkFixWinding(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Polygon{
			Line{Point{X: -2.109375, Y: 11.178401873711785}, Point{X: -16.875, Y: -43.06888777416961}, Point{X: 62.57812500000001, Y: -43.580390855607845}, Point{X: 81.5625, Y: 8.407168163601076}},
			Line{Point{X: 7.3828125, Y: -23.241346102386135}, Point{X: 28.4765625, Y: -8.05922962720018}, Point{X: 55.1953125, Y: -11.178401873711772}, Point{X: 22.148437499999996, Y: -33.137551192346145}},
			Line{Point{X: 25.48828125, Y: -18.312810846425432}, Point{X: 33.22265625, Y: -16.720385051693988}, Point{X: 34.013671875, Y: -21.207458730482642}, Point{X: 23.466796875, Y: -24.766784522874428}},
			Line{Point{X: 27.5537109375, Y: -12.618897304044012}, Point{X: 29.02587890625, Y: -12.146745814539685}, Point{X: 29.377441406249996, Y: -14.604847155053898}, Point{X: 26.3671875, Y: -15.855673509998681}},
			Line{Point{X: 27.0703125, Y: -20.3034175184893}, Point{X: 27.509765625, Y: -21.616579336740593}, Point{X: 31.113281249999996, Y: -19.559790136497398}}}.FixWinding()
	}
}

func TestPolygonValidTopology(t *testing.T) {
	p := Polygon{Line{{3, 4}, {2, 9}, {1, 4}}}
	assert.True(t, p.ValidTopology())

	p = Polygon{Line{{3, 4}, {2, 9}, {1, 4}, {1, 5}}}
	assert.False(t, p.ValidTopology())
}

func TestPolygonClipBBoxShortCircuit(t *testing.T) {
	t.Run("completely inside bbox", func(t *testing.T) {
		p := Polygon{Line{{1, 1}, {2, 1}, {2, 2}, {1, 2}}}
		bbox := BBox{SW: Point{0, 0}, NE: Point{3, 3}}

		assert.Equal(t,
			[]Geom{MustNewGeom(Polygon{Line{
				{1, 1}, {2, 1}, {2, 2}, {1, 2},
			}})},
			p.ClipToBBox(bbox),
		)
	})

	t.Run("fit to bbox", func(t *testing.T) {
		p := Polygon{Line{{0, 0}, {3, 0}, {3, 3}, {0, 3}}}
		bbox := BBox{SW: Point{1, 1}, NE: Point{2, 2}}

		assert.Equal(t,
			[]Geom{MustNewGeom(Polygon{Line{
				{1, 1}, {2, 1}, {2, 2}, {1, 2},
			}})},
			p.ClipToBBox(bbox),
		)
	})

	t.Run("no speedup", func(t *testing.T) {
		p := Polygon{Line{{0, 0}, {3, 0}, {0, 3}}}
		bbox := BBox{SW: Point{1, 1}, NE: Point{2, 2}}

		assert.Equal(t,
			[]Geom{MustNewGeom(Polygon{Line{
				{1, 1}, {1, 2}, {2, 1},
			}})},
			p.ClipToBBox(bbox),
		)
	})
}
