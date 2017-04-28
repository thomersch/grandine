package spatial

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalWKB(t *testing.T) {
	f := Feature{}
	buf, err := f.MarshalWKB()
	assert.Nil(t, err)

	fmt.Printf("%v\n", buf)
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
	fmt.Printf("%s\n", buf)
}
