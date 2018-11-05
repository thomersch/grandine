package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/thomersch/grandine/lib/spatial"
)

const defaultChunkSize = 100000

type csvReader interface {
	Read() ([]string, error)
}

type Chunks struct {
	chunkSize   int
	lineDecoder func(csvReader) (spatial.Feature, error)
	csvRdr      csvReader

	endReached bool
}

func (c *Chunks) Next() bool {
	return !c.endReached
}

func (c *Chunks) Scan(fc *spatial.FeatureCollection) error {
	if c.chunkSize == 0 {
		c.chunkSize = defaultChunkSize
	}
	var (
		fts  = make([]spatial.Feature, c.chunkSize)
		err  error
		read int
	)
	for i := range fts {
		fts[i], err = c.lineDecoder(c.csvRdr)
		if err == io.EOF {
			c.endReached = true
			break
		}
		if err != nil {
			return err
		}
		read++
	}
	fc.Features = append(fc.Features, fts[:read]...)
	return nil
}

type Codec struct {
	LatCol, LonCol int
	Delim          rune

	keys []string
}

func (c *Codec) decodeLine(csvr csvReader) (spatial.Feature, error) {
	var (
		ft = spatial.NewFeature()
		pt spatial.Point
	)

	record, err := csvr.Read()
	if err != nil {
		return ft, err
	}

	if c.LonCol >= len(record) || c.LatCol >= len(record) {
		return ft, fmt.Errorf("there are not enough columns in: '%v'", record)
	}

	pt.X, err = strconv.ParseFloat(record[c.LonCol], 64)
	if err != nil {
		return ft, err
	}
	pt.Y, err = strconv.ParseFloat(record[c.LatCol], 64)
	if err != nil {
		return ft, err
	}
	ft.Geometry = spatial.MustNewGeom(pt)
	for i, val := range record {
		if i >= len(c.keys) {
			// there are more value in this line than header keys
			continue
		}
		ft.Props[c.keys[i]] = val
	}
	return ft, nil
}

func (c *Codec) newReader(r io.Reader) csvReader {
	csvRdr := csv.NewReader(r)
	if c.Delim == 0 {
		csvRdr.Comma = '	'
	} else {
		csvRdr.Comma = c.Delim
	}
	csvRdr.LazyQuotes = false
	return csvRdr
}

func (c *Codec) Decode(r io.Reader, fs *spatial.FeatureCollection) error {
	csvRdr := c.newReader(r)

	var err error
	c.keys, err = csvRdr.Read()
	if err != nil {
		return err
	}

	for {
		ft, err := c.decodeLine(csvRdr)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fs.Features = append(fs.Features, ft)
	}
	return nil
}

func (c *Codec) ChunkedDecode(r io.Reader) (spatial.Chunks, error) {
	csvRdr := c.newReader(r)

	var err error
	c.keys, err = csvRdr.Read()
	if err != nil {
		return nil, err
	}

	return &Chunks{
		chunkSize:   defaultChunkSize,
		lineDecoder: c.decodeLine,
		csvRdr:      csvRdr,
	}, nil
}

func (c *Codec) Extensions() []string {
	return []string{"csv", "txt"}
}
