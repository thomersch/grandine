// +build !golangclip

package spatial

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeosConversion(t *testing.T) {
	ply := Polygon{{{0, 1}, {2.5, 1}, {3, 2}, {2.5, 3}, {1, 3}}}
	gs := ply.geos()
	assert.Equal(t, ply, geosToPolygon(gs))
}

func TestGeosSelfIntersect(t *testing.T) {
	f, err := os.Open("testfiles/self_intersect.geojson")
	assert.Nil(t, err)
	defer f.Close()

	fc := NewFeatureCollection()
	err = json.NewDecoder(f).Decode(&fc)
	assert.Nil(t, err)

	res := fc.Features[0].Geometry.MustPolygon().clipToBBox(BBox{Point{0, 0}, Point{2000, 2000}})
	assert.Len(t, res, 2) // TODO: some more testing might be useful
}
