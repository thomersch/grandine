package csv

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/thomersch/grandine/lib/spaten/fileformat"
	"github.com/thomersch/grandine/lib/spatial"
)

type TagMapping struct {
	Name string
	Type fileformat.Tag_ValueType
}

type Codec struct {
	LatCol, LonCol int
	Delim          rune
}

func (c *Codec) Decode(r io.Reader, fs *spatial.FeatureCollection) error {
	csvrdr := csv.NewReader(r)
	if c.Delim == 0 {
		csvrdr.Comma = '	'
	} else {
		csvrdr.Comma = c.Delim
	}
	csvrdr.LazyQuotes = false

	keys, err := csvrdr.Read()
	if err != nil {
		return err
	}

	for {
		record, err := csvrdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var (
			ft = spatial.NewFeature()
			pt spatial.Point
		)
		pt.X, err = strconv.ParseFloat(record[c.LonCol], 64)
		if err != nil {
			return err
		}
		pt.Y, err = strconv.ParseFloat(record[c.LatCol], 64)
		if err != nil {
			return err
		}
		ft.Geometry = spatial.MustNewGeom(pt)
		for i, val := range record {
			ft.Props[keys[i]] = val
		}
		fs.Features = append(fs.Features, ft)
	}
	return nil
}

func (c *Codec) Extensions() []string {
	return []string{"csv", "txt"}
}
