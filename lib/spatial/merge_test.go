package spatial

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	t.Run("different key-value", func(t *testing.T) {
		props1 := map[string]interface{}{
			"1": 2,
		}
		props2 := map[string]interface{}{
			"1": 3.1,
		}
		feats := []Feature{
			{
				Props:    props1,
				Geometry: MustNewGeom(Line{{1, 2}, {3, 4}}),
			},
			{
				Props:    props2,
				Geometry: MustNewGeom(Line{{3, 4}, {5, 6}}),
			},
		}
		assert.Equal(t, feats, MergeFeatures(feats))
	})

	t.Run("continuous", func(t *testing.T) {
		props := map[string]interface{}{
			"a": 1,
			"b": "foo",
			"c": 1.234,
		}
		feat1 := Feature{
			Props:    props,
			Geometry: MustNewGeom(Line{{1, 0}, {1, 1}, {2, 3}, {5, 6}}),
		}
		feat2 := Feature{
			Props:    props,
			Geometry: MustNewGeom(Line{{5, 6}, {7, 8}, {6, 6}, {4, 5}}),
		}
		assert.Equal(t, []Feature{
			{
				Props:    props,
				Geometry: MustNewGeom(Line{{1, 0}, {1, 1}, {2, 3}, {5, 6}, {7, 8}, {6, 6}, {4, 5}}),
			},
		}, MergeFeatures([]Feature{feat1, feat2}))
	})
}

func TestMergeFoo(t *testing.T) {
	f, err := os.Open("testfiles/mergable_lines.geojson")
	assert.Nil(t, err)

	var fcoll = NewFeatureCollection()
	err = json.NewDecoder(f).Decode(fcoll)
	assert.Nil(t, err)

	fcoll.Features = searchAndMerge(fcoll.Features)

	assert.Len(t, fcoll.Features, 1)
	assert.True(t, len(fcoll.Features[0].Geometry.MustLineString()) > 7)
}

func TestBuckets(t *testing.T) {
	f1 := Feature{Props: map[string]interface{}{"1": 2}}
	f2 := Feature{Props: map[string]interface{}{"1": 2}}
	f3 := Feature{Props: map[string]interface{}{"1": 3}}

	bucks := tagBuckets([]Feature{f1, f2, f3})
	assert.Len(t, bucks, 2)
}

func TestIgnoreList(t *testing.T) {
	il := ignoreList{1, 5, 9}
	il.Add(3)
	assert.Equal(t, ignoreList{1, 3, 5, 9}, il)
	il.Add(19)
	assert.Equal(t, ignoreList{1, 3, 5, 9, 19}, il)
	assert.True(t, il.Has(1))
	assert.False(t, il.Has(4))
}
