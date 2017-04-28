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
}
