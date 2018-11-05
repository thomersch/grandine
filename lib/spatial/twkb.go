package spatial

import (
	"encoding/binary"
	"io"
	"math"
)

func twkbWriteHeader(w io.Writer, gt GeomType, precision int) error {
	var buf = make([]byte, 2)
	buf[0] = byte(precision<<4) ^ byte(gt)
	_, err := w.Write(buf)
	return err
}

func twkbWritePoint(w io.Writer, p Point, previous Point, precision int) error {
	var (
		xi  = int(p.X * math.Pow10(precision))
		yi  = int(p.Y * math.Pow10(precision))
		xpi = int(previous.X * math.Pow10(precision))
		ypi = int(previous.Y * math.Pow10(precision))

		buf = make([]byte, 20) // up to 10 bytes per varint
	)

	dx := int64(xi - xpi)
	dy := int64(yi - ypi)
	bwx := binary.PutVarint(buf, dx)
	bwy := binary.PutVarint(buf[bwx:], dy)

	_, err := w.Write(buf[:bwx+bwy])
	return err
}

type combinedReader interface {
	io.Reader
	io.ByteReader
}

type wrappedReader struct {
	io.Reader
}

func (wr *wrappedReader) ReadByte() (c byte, err error) {
	var b = make([]byte, 1)
	n, err := wr.Read(b)
	if err != nil {
		return b[0], err
	}
	if n != 1 {
		return b[0], io.EOF
	}
	return b[0], nil
}

type twkbHeader struct {
	typ       GeomType
	precision int
	// metadata attributes
	bbox, size, idList, extendedPrecision, emptyGeom bool
}

func unzigzag(nVal int) int {
	if (nVal & 1) == 0 {
		return nVal >> 1
	}
	return -(nVal >> 1) - 1
}

func twkbReadHeader(r io.Reader) (twkbHeader, error) {
	var (
		// BIT   USAGE
		// 1-4   type
		// 5-8   precision
		// 9	 bbox
		// 10    size attribute
		// 11    id list
		// 12    extended precision
		// 13    empty geom
		// 14-16 unused
		by = make([]byte, 2)
		hd twkbHeader
	)
	_, err := r.Read(by)
	hd.typ = GeomType(by[0] & 15)
	hd.precision = int(by[0] >> 4)
	hd.bbox = int(by[1])&1 == 1
	hd.size = int(by[1])&2 == 2
	hd.idList = int(by[1])&4 == 4
	hd.extendedPrecision = int(by[1])&8 == 8
	hd.emptyGeom = int(by[1])&16 == 16
	return hd, err
}

func twkbReadPoint(r io.Reader, previous Point, precision int) (Point, error) {
	wr, ok := r.(io.ByteReader)
	if !ok {
		wr = &wrappedReader{r}
	}
	var pt Point
	xe, err := binary.ReadVarint(wr)
	if err != nil {
		return pt, err
	}
	ye, err := binary.ReadVarint(wr)
	if err != nil {
		return pt, err
	}
	xΔ := float64(xe) / math.Pow10(precision)
	yΔ := float64(ye) / math.Pow10(precision)

	pt.X = xΔ + previous.X
	pt.Y = yΔ + previous.Y
	return pt, nil
}

func twkbReadLineString(r io.Reader, precision int) ([]Point, error) {
	wr, ok := r.(combinedReader)
	if !ok {
		wr = &wrappedReader{r}
	}
	npoints, err := binary.ReadUvarint(wr)
	if err != nil {
		return nil, err
	}

	var (
		ls     = make([]Point, npoints)
		lastPt Point
	)
	for i := 0; i < int(npoints); i++ {
		lastPt, err = twkbReadPoint(wr, lastPt, precision)
		if err != nil {
			return ls, err
		}
		ls[i] = lastPt
	}
	return ls, nil
}

func twkbWriteLineString(w io.Writer, ls []Point, precision int) error {
	buf := make([]byte, 10)
	bWritten := binary.PutUvarint(buf, uint64(len(ls)))
	_, err := w.Write(buf[:bWritten-1])
	if err != nil {
		return err
	}
	var previous Point
	for _, pt := range ls {
		if err = twkbWritePoint(w, pt, previous, precision); err != nil {
			return err
		}
	}
	return nil
}
