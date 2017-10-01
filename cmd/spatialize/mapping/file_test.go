package mapping

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMapping(t *testing.T) {
	f, err := os.Open("mapping.yml")
	assert.Nil(t, err)

	_, err = ParseMapping(f)
	assert.Nil(t, err)
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
