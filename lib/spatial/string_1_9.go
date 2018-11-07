// +build !go1.10

package spatial

import "fmt"

func (p Polygon) string() string {
	s := "("
	for _, line := range p {
		s += fmt.Sprintf("%v, ", line)
	}
	return s[:len(s)-2] + ")"
}

func (l Line) string() string {
	s := ""
	for _, point := range l {
		s += fmt.Sprintf("%v, ", point)
	}
	return s[:len(s)-2]
}
