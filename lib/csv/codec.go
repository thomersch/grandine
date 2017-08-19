package csv

import (
	"encoding/csv"
	"io"
	"log"
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
	ColPropMap     map[int]TagMapping
}

func (c *Codec) Decode(r io.Reader, fs *spatial.FeatureCollection) error {
	csvrdr := csv.NewReader(r)
	csvrdr.Comma = '	'
	csvrdr.LazyQuotes = false

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
		pt[0], err = strconv.ParseFloat(record[c.LonCol], 64)
		if err != nil {
			return err
		}
		pt[1], err = strconv.ParseFloat(record[c.LatCol], 64)
		if err != nil {
			return err
		}
		ft.Geometry = spatial.MustNewGeom(pt)
		for i, mapping := range c.ColPropMap {
			var (
				v   interface{}
				err error
			)
			switch mapping.Type {
			case fileformat.Tag_STRING:
				v = record[i]
			case fileformat.Tag_INT:
				v, err = strconv.Atoi(record[i])
			}
			if err != nil {
				log.Println("could not parse '%s': %v", record[i], err)
			}
			ft.Props[mapping.Name] = v
		}
		fs.Features = append(fs.Features, ft)
	}
	return nil
}

func (c *Codec) Extensions() []string {
	return []string{"csv", "txt"}
}
