package mvt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomersch/grandine/lib/spatial"
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
			expectedResult: []uint32{9, 2, 2},
		},
		{
			geom: []interface{}{
				spatial.Point{25, 17},
			},
			expectedResult: []uint32{9, 50, 34},
		},
	}

	for n, tc := range tcs {
		t.Run(fmt.Sprintf("%v", n), func(t *testing.T) {
			var fts []spatial.Feature
			for _, g := range tc.geom {
				geom, err := spatial.NewGeom(g)
				assert.Nil(t, err)
				fts = append(fts, spatial.Feature{Geometry: geom})
			}
			res := encodeGeometry(fts)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
