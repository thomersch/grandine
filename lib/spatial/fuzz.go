// +build gofuzz

package spatial

import "bytes"

func Fuzz(data []byte) int {
	var g Geom

	err := g.UnmarshalWKB(bytes.NewReader(data))
	if err != nil {
		return 0
	}
	return 1
}
