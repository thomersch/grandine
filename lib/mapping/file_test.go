package mapping

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMapping(t *testing.T) {
	f, err := os.Open("mapping.yml")
	assert.Nil(t, err)

	conds, err := ParseMapping(f)
	assert.Nil(t, err)

	srcKV := map[string]interface{}{
		"building": "yes",
	}
	assert.True(t, conds[1].Matches(srcKV))
	assert.Equal(t, map[string]interface{}{
		"@layer":    "building",
		"@zoom:min": 14,
	}, conds[1].Map(srcKV))

	srcKV = map[string]interface{}{
		"highway": "primary",
	}
	assert.True(t, conds[0].Matches(srcKV))
	assert.Equal(t, map[string]interface{}{
		"@layer": "transportation",
		"class":  "primary",
	}, conds[0].Map(srcKV))

	srcKV = map[string]interface{}{
		"railway":  "rail",
		"maxspeed": "300",
	}
	assert.True(t, conds[2].Matches(srcKV))
	assert.Equal(t, map[string]interface{}{
		"@layer":   "transportation",
		"class":    "railway",
		"maxspeed": 300,
	}, conds[2].Map(srcKV))
}
