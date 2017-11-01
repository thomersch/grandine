package spatial

import (
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
