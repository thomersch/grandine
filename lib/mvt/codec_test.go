package mvt

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"

	"github.com/stretchr/testify/assert"
)

func TestEncodeGeometry(t *testing.T) {
	tcs := []struct {
		geom           []interface{}
		expectedResult []uint32
	}{
		{
			geom: []interface{}{
				spatial.Point{1, 1},
			},
			// TODO: validate coordinates
			expectedResult: []uint32{9, 44, 8148},
		},
		{
			geom: []interface{}{
				spatial.Point{25, 17},
			},
			// TODO: validate coordinates
			expectedResult: []uint32{9, 1136, 7408},
		},
	}

	for n, tc := range tcs {
		t.Run(fmt.Sprintf("%v", n), func(t *testing.T) {
			var geoms []spatial.Geom
			for _, g := range tc.geom {
				geom, err := spatial.NewGeom(g)
				assert.Nil(t, err)
				geoms = append(geoms, geom)
			}
			res, err := encodeGeometry(geoms, tile.ID{X: 1, Y: 0, Z: 1})
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestEncodeTile(t *testing.T) {
	var features []spatial.Feature
	geoms := []interface{}{
		// Point
		spatial.Point{45, 45},
		spatial.Point{50, 47},
		spatial.Point{100, 40},
		spatial.Point{179, 40},
		// LineString
		spatial.Line{
			spatial.Point{
				-1.0546875,
				55.97379820507658,
			},
			spatial.Point{
				14.765625,
				44.08758502824516,
			},
			spatial.Point{
				39.7265625,
				67.7427590666639,
			},
			spatial.Point{
				16.875,
				67.06743335108297,
			},
			spatial.Point{
				16.171875,
				58.07787626787517,
			},
		},
		spatial.Polygon{spatial.Line{
			spatial.Point{
				2.8125,
				54.77534585936447,
			},
			spatial.Point{
				1.23046875,
				47.87214396888731,
			},
			spatial.Point{
				7.207031249999999,
				37.020098201368114,
			},
			spatial.Point{
				21.26953125,
				40.97989806962013,
			},
			spatial.Point{
				29.8828125,
				48.69096039092549,
			},
			spatial.Point{
				31.113281249999996,
				53.12040528310657,
			},
			spatial.Point{
				23.90625,
				60.413852350464914,
			},
			spatial.Point{
				10.01953125,
				60.84491057364915,
			},
			spatial.Point{
				2.8125,
				54.77534585936447,
			},
		}},
	}

	for _, geom := range geoms {
		g, err := spatial.NewGeom(geom)
		assert.Nil(t, err)
		features = append(features, spatial.Feature{Geometry: g})
	}

	features[0].Props = map[string]interface{}{
		"highway": "primary",
		"oneway":  1,
	}
	features[1].Props = map[string]interface{}{
		"highway": "secondary",
		"oneway":  -1,
	}
	features[2].Props = map[string]interface{}{
		"ignorance": "strength",
	}

	layers := map[string][]spatial.Feature{
		"main": features,
	}

	buf, err := EncodeTile(layers, tile.ID{X: 1, Y: 0, Z: 1})
	assert.Nil(t, err)

	var b bytes.Buffer
	b.Write(buf)
}

func TestEncodeLine(t *testing.T) {
	var cur [2]int
	ln := encodeLine(spatial.Line{{0, 1}, {3, 4}, {10, 1}}, &cur, 10, 100000, 100000, 0, 0)
	assert.Equal(t, []uint32{9, 0, 1, 18, 666, 7, 1560, 8}, ln)
}
