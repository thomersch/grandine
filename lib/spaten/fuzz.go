// +build gofuzz

package spaten

import (
	"bytes"

	"github.com/thomersch/grandine/lib/spatial"
)

func Fuzz(data []byte) int {
	var (
		c  Codec
		fc = spatial.NewFeatureCollection()
	)

	err := c.Decode(bytes.NewReader(data), fc)
	if err != nil {
		return 0
	}
	return 1
}
