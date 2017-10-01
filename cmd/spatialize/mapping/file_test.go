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

	srcKV := map[string]string{
		"building": "yes",
	}
	assert.True(t, conds[1].Matches(srcKV))
	assert.Equal(t, map[string]interface{}{
		"@layer":    "building",
		"@zoom:min": 14,
	}, conds[1].Map(srcKV))

	srcKV = map[string]string{
		"highway": "primary",
	}
	assert.True(t, conds[0].Matches(srcKV))
	assert.Equal(t, map[string]interface{}{
		"@layer": "transportation",
		"class":  "primary",
	}, conds[0].Map(srcKV))
}

// o, err := yaml.Marshal(fileMappings{
// 	{
// 		Src: fileMapKV{
// 			Key:   "sk",
// 			Value: "sv",
// 		},
// 		Dest: []fileMapKV{
// 			fileMapKV{Key: "dk1", Value: "dv"},
// 			fileMapKV{Key: "dk2", Value: "dv"},
// 		},
// 	},
// 	{
// 		Src: fileMapKV{
// 			Key:   "sk",
// 			Value: "sv",
// 		},
// 		Dest: []fileMapKV{
// 			fileMapKV{Key: "dk1", Value: "dv"},
// 			fileMapKV{Key: "dk2", Value: "dv"},
// 		},
// 	},
// })
// if err != nil {
// 	return nil, err
// }
// fmt.Printf("%s\n", o)
