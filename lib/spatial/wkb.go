package spatial

import (
	"encoding/binary"
	"fmt"
	"io"
)

func wkbWritePoint(w io.Writer, p Point) {
	fmt.Println(p)
	binary.Write(w, endianness, float64(p.X()))
	binary.Write(w, endianness, float64(p.Y()))
}

func wkbWriteLineString(w io.Writer, ls []Point) {
	binary.Write(w, endianness, uint32(len(ls)))
	for _, pt := range ls {
		binary.Write(w, endianness, pt.X())
		binary.Write(w, endianness, pt.Y())
	}
}

func wkbWritePolygon(w io.Writer, poly [][]Point) {

}
