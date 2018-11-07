// +build go1.10

package spatial

import (
	"fmt"
	"strings"
)

func stringPolygon(p Polygon) string {
	var b strings.Builder
	b.WriteByte('(')
	for pos, line := range p {
		fmt.Fprintf(&b, "%v", line)
		if pos != len(p)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteByte(')')
	return b.String()
}

func stringLine(l Line) string {
	var b strings.Builder
	for pos, point := range l {
		fmt.Fprintf(&b, "%v", point)
		if pos != len(l)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}
