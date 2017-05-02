package spatial

import (
	"io"
	"math"
)

const wkbRawPointSize = 16

func wkbWritePoint(w io.Writer, p Point) error {
	var (
		err error
		buf = make([]byte, wkbRawPointSize)
	)
	endianness.PutUint64(buf[:8], math.Float64bits(p.X()))
	endianness.PutUint64(buf[8:16], math.Float64bits(p.Y()))
	_, err = w.Write(buf)
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

func wkbReadHeader(r io.Reader) (GeomType, error) {
	var buf = make([]byte, 4)
	_, err := r.Read(buf)
	gt := endianness.Uint32(buf)
	return GeomType(gt), err
}

// TODO: evaluate returning Geom instead of Point
func wkbReadPoint(r io.Reader) (p Point, err error) {
	var buf = make([]byte, wkbRawPointSize)
	n, err := r.Read(buf)
	if n != wkbRawPointSize {
		return p, io.EOF
	}
	if err != nil {
		return
	}
	p[0] = math.Float64frombits(endianness.Uint64(buf[:8]))
	p[1] = math.Float64frombits(endianness.Uint64(buf[8:16]))
	return
}

// TODO: evaluate returning Geom instead of Point
func wkbReadLineString(r io.Reader) ([]Point, error) {
	var buf = make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	nop := endianness.Uint32(buf)

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
	var buf = make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	nor := endianness.Uint32(buf)

	var rings = make([][]Point, nor)
	for i := 0; i < int(nor); i++ {
		rings[i], err = wkbReadLineString(r)
		if err != nil {
			return rings, err
		}
	}
	return rings, err
}
