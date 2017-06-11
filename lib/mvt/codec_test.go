package mvt

import (
	"fmt"
	"os"
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

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
			expectedResult: []uint32{9, 182, 198},
		},
		{
			geom: []interface{}{
				spatial.Point{25, 17},
			},
			// TODO: validate coordinates
			expectedResult: []uint32{9, 4550, 3398},
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
			res, err := encodeGeometry(geoms, TileID{X: 4, Y: 4, Z: 3})
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestEncodeTile(t *testing.T) {
	var features []spatial.Feature
	geoms := []interface{}{
		spatial.Point{10, 10},
		spatial.Point{20, 10},
		spatial.Point{20, 30},
	}

	for _, geom := range geoms {
		g, err := spatial.NewGeom(geom)
		assert.Nil(t, err)
		features = append(features, spatial.Feature{Geometry: g})
	}

	layers := map[string][]spatial.Feature{
		"main": features,
	}

	buf, err := EncodeTile(layers, TileID{X: 0, Y: 0, Z: 1})
	assert.Nil(t, err)
	f, err := os.Create("testtile.mvt")
	assert.Nil(t, err)
	defer f.Close()
	f.Write(buf)
}
