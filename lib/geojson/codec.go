package geojson

import (
	"encoding/json"
	"errors"
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
	if srid, ok := ogcSRID[gjfc.CRS.Properties.Name]; ok {
		if len(fc.SRID) == 0 {
			fc.SRID = srid
		}
		if len(srid) != 0 && fc.SRID != srid {
			return errors.New("incompatible projections: ")
		}
	}
	return nil
}

func (c *Codec) Encode(w io.Writer, fc *spatial.FeatureCollection) error {
	geojsonFC := featureColl{
		Type:     "FeatureCollection",
		Features: fc.Features,
	}
	if len(fc.SRID) != 0 {
		geojsonFC.CRS.Properties.Name, _ = sridOGC[fc.SRID]
		geojsonFC.CRS.Type = "name"
	}
	return json.NewEncoder(w).Encode(&geojsonFC)
}

func (c *Codec) Extensions() []string {
	return []string{"geojson", "json"}
}

type featureColl struct {
	Type string `json:"type"`
	CRS  struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
	Features []spatial.Feature `json:"features"`
}
