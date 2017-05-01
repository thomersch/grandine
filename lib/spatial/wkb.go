package spatial

import (
	"encoding/binary"
	"io"
	"math"
)

func wkbWritePoint(w io.Writer, p Point) error {
	var (
		err error
		x   = make([]byte, 8)
		y   = make([]byte, 8)
	)
	endianness.PutUint64(x, math.Float64bits(p.X()))
	endianness.PutUint64(y, math.Float64bits(p.Y()))
	_, err = w.Write(append(x, y...))
	if err != nil {
		return err
	}
	return nil
}

func wkbWriteLineString(w io.Writer, ls []Point) error {
	// write number of points
	var ln = make([]byte, 4)
	endianness.PutUint32(ln, uint32(len(ls)))
	_, err := w.Write(ln)
	if err != nil {
		return err
	}

	for _, pt := range ls {
		err = wkbWritePoint(w, pt)
		if err != nil {
			return err
		}
	}
	return nil
}

func wkbWritePolygon(w io.Writer, poly [][]Point) error {
	// write number of rings
	var lnr = make([]byte, 4)
	endianness.PutUint32(lnr, uint32(len(poly)))
	_, err := w.Write(lnr)
	if err != nil {
		return err
	}

	for _, ring := range poly {
		err = wkbWriteLineString(w, ring)
		if err != nil {
			return err
		}
	}
	return nil
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
