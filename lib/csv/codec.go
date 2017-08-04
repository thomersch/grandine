package csv

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/thomersch/grandine/lib/spatial"
)

type reader struct {
	LatCol, LonCol int
	ColPropMap     map[int]string
}

func (rdr *reader) Decode(r io.Reader, fs *spatial.FeatureCollection) error {
	csvrdr := csv.NewReader(r)
	csvrdr.Comma = '	'

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
		pt[0], err = strconv.ParseFloat(record[rdr.LonCol], 64)
		if err != nil {
			return err
		}
		pt[1], err = strconv.ParseFloat(record[rdr.LatCol], 64)
		if err != nil {
			return err
		}
		ft.Geometry = spatial.MustNewGeom(pt)
		for i, keyName := range rdr.ColPropMap {
			ft.Props[keyName] = record[i]
		}
		fs.Features = append(fs.Features, ft)
	}
	return nil
}
