package geojson

import (
	"encoding/json"
	"io"

	"github.com/thomersch/grandine/lib/spatial"
)

type Codec struct{}

func (c *Codec) Decode(r io.Reader, fc *spatial.FeatureCollection) error {
	var ffc spatial.FeatureCollection
	err := json.NewDecoder(r).Decode(&ffc)
	if err != nil {
		return err
	}
	fc.Features = append(fc.Features, ffc.Features...)
	return nil
}

func (c *Codec) Encode(w io.Writer, fc *spatial.FeatureCollection) error {
	return json.NewEncoder(w).Encode(&fc)
}

func (c *Codec) Extensions() []string {
	return []string{"geojson", "json"}
}
