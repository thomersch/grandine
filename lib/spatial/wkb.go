package spatial

import (
	"errors"
	"fmt"
	"io"
	"math"
)

type WKBable interface {
	MarshalWKB() ([]byte, error)
}

type WKBableWithProps interface {
	WKBable
	PropertyRetriever
}

const (
	wkbRawPointSize = 16
	wkbLittleEndian = 1
)

func GeomFromWKB(r io.Reader) (Geom, error) {
	var (
		g             Geom
		wkbEndianness = make([]byte, 1)
	)
	_, err := r.Read(wkbEndianness)
	if err != nil {
		return g, err
	}
	if wkbEndianness[0] != wkbLittleEndian {
		return g, errors.New("only little endian is supported")
	}

	g.typ, err = wkbReadHeader(r)
	if err != nil {
		return g, err
	}
	switch g.typ {
	case GeomTypePoint:
		g.g, err = wkbReadPoint(r)
	case GeomTypeLineString:
		g.g, err = wkbReadLineString(r)
	case GeomTypePolygon:
		g.g, err = wkbReadPolygon(r)
	default:
		return g, fmt.Errorf("unsupported GeomType: %v", g.typ)
	}
	return g, err
}

func wkbReadHeader(r io.Reader) (GeomType, error) {
	var buf = make([]byte, 4)
	_, err := r.Read(buf)
	gt := endianness.Uint32(buf)
	return GeomType(gt), err
}

func wkbWritePoint(w io.Writer, p Point) error {
	var (
		err error
		buf = make([]byte, wkbRawPointSize)
	)
	endianness.PutUint64(buf[:8], math.Float64bits(p.X))
	endianness.PutUint64(buf[8:16], math.Float64bits(p.Y))
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

func wkbWritePolygon(w io.Writer, poly Polygon) error {
	// write number of rings
	var lnr = make([]byte, 4)
	endianness.PutUint32(lnr, uint32(len(poly)))
	_, err := w.Write(lnr)
	if err != nil {
		return err
	}

	for _, ring := range poly {
		err = wkbWriteLineString(w, append(ring, ring[0])) // wkb closes rings with the first element, the internal implementation doesn't
		if err != nil {
			return err
		}
	}
	return nil
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
	p.X = math.Float64frombits(endianness.Uint64(buf[:8]))
	p.Y = math.Float64frombits(endianness.Uint64(buf[8:16]))
	return
}

// TODO: evaluate returning Geom instead of Point
func wkbReadLineString(r io.Reader) (Line, error) {
	var buf = make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	nop := endianness.Uint32(buf)
	if nop == 0 {
		return nil, errors.New("a linestring needs to have at least one point")
	}

	var ls = make(Line, nop)
	for i := 0; i < int(nop); i++ {
		ls[i], err = wkbReadPoint(r)
		if err != nil {
			return ls, err
		}
	}
	return ls, nil
}

func wkbReadPolygon(r io.Reader) (Polygon, error) {
	var buf = make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	nor := endianness.Uint32(buf)
	if nor == 0 {
		return nil, errors.New("a polygon needs to have at least one ring")
	}

	var rings = make(Polygon, nor)
	for i := 0; i < int(nor); i++ {
		rings[i], err = wkbReadLineString(r)
		if err != nil {
			return rings, err
		}
		rings[i] = rings[i][:len(rings[i])-1] // wkb closes rings with the first element, the internal implementation doesn't
	}
	return rings, err
}
