package csv

import (
	"os"
	"testing"

	"github.com/thomersch/grandine/lib/spatial"

	"github.com/stretchr/testify/assert"
)

func TestCSVDecode(t *testing.T) {
	f, err := os.Open("testfiles/gn_excerpt.csv")
	assert.Nil(t, err)

	csvr := reader{
		LatCol: 4,
		LonCol: 5,
		ColPropMap: map[int]string{
			1: "name",
		},
	}
	fcoll := spatial.FeatureCollection{}
	err = csvr.Decode(f, &fcoll)
	assert.Nil(t, err)
	assert.Equal(t, "les Escaldes", fcoll.Features[0].Props["name"])
	pt, err := fcoll.Features[0].Geometry.Point()
	assert.Nil(t, err)
	assert.Equal(t, 1.53414, pt.X())
	assert.Equal(t, 42.50729, pt.Y())
}
