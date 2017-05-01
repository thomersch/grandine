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

// TODO: evaluate returning Geom instead of Point
func wkbReadPoint(r io.Reader) (p Point, err error) {
	err = binary.Read(r, endianness, &p[0])
	if err != nil {
		return
	}
	err = binary.Read(r, endianness, &p[1])
	return
}

// TODO: evaluate returning Geom instead of Point
func wkbReadLineString(r io.Reader) ([]Point, error) {
	var nop uint32 // number of points
	err := binary.Read(r, endianness, &nop)
	if err != nil {
		return nil, err
	}

	var ls = make([]Point, nop)
	for i := 0; i < int(nop); i++ {
		ls[i], err = wkbReadPoint(r)
		if err != nil {
			return ls, err
		}
	}
	return ls, nil
}

func wkbReadPolygon(r io.Reader) ([][]Point, error) {
	var nor uint32 // number of rings
	err := binary.Read(r, endianness, &nor)
	if err != nil {
		return nil, err
	}

	var rings = make([][]Point, nor)
	for i := 0; i < int(nor); i++ {
		rings[i], err = wkbReadLineString(r)
		if err != nil {
			return rings, err
		}
	}
	return rings, err
}
