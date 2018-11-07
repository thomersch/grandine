// +build !go1.10

package spatial

import "fmt"

func stringPolygon(p Polygon) string {
	s := "("
	for _, line := range p {
		s += fmt.Sprintf("%v, ", line)
	}
	return s[:len(s)-2] + ")"
}

func stringLine(l Line) string {
	s := ""
	for _, point := range l {
		s += fmt.Sprintf("%v, ", point)
	}
	return s[:len(s)-2]
}
