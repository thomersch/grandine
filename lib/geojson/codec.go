package geojson

import (
	"encoding/json"
	"io"

	"github.com/thomersch/grandine/lib/spatial"
)

type Codec struct{}

func (c *Codec) Decode(r io.Reader, fc *spatial.FeatureCollection) error {
	var gjfc featureColl
	err := json.NewDecoder(r).Decode(&gjfc)
	if err != nil {
		return err
	}
	fc.Features = append(fc.Features, gjfc.Features...)
	// TODO: set SRID
	return nil
}

func (c *Codec) Encode(w io.Writer, fc *spatial.FeatureCollection) error {
	geojsonFC := featureColl{
		Type: "FeatureCollection",
		// TODO: set CRS
		Features: fc.Features,
	}
	return json.NewEncoder(w).Encode(&geojsonFC)
}

func (c *Codec) Extensions() []string {
	return []string{"geojson", "json"}
}

type featureColl struct {
	Type string `json:"type"`
	CRS  struct {
		Type       string
		Properties struct {
			Name string
		}
	} `json:"crs"`
	Features []spatial.Feature
}
