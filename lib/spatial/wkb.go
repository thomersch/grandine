package spatial

import (
	"encoding/binary"
	"io"
)

func wkbWritePoint(w io.Writer, p Point) {
	binary.Write(w, endianness, float64(p.X()))
	binary.Write(w, endianness, float64(p.Y()))
}

func wkbWriteLineString(w io.Writer, ls []Point) {
	// write number of points
	binary.Write(w, endianness, uint32(len(ls)))
	for _, pt := range ls {
		binary.Write(w, endianness, pt.X())
		binary.Write(w, endianness, pt.Y())
	}
}

func wkbWritePolygon(w io.Writer, poly [][]Point) {
	// write number of rings
	binary.Write(w, endianness, uint32(len(poly)))
	for _, ring := range poly {
		wkbWriteLineString(w, ring)
	}
}
