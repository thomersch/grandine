package geojson

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

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
	Features featList `json:"features"`
}

type featureProto struct {
	Geometry struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	Properties map[string]interface{} `json:"properties"`
}

type featList []spatial.Feature

func (fl *featList) UnmarshalJSON(buf []byte) error {
	var fts []featureProto
	err := json.Unmarshal(buf, &fts)
	if err != nil {
		return err
	}

	*fl = make([]spatial.Feature, 0, len(fts))
	for _, inft := range fts {
		if singularType(inft.Geometry.Type) {
			var ft spatial.Feature
			ft.Props = inft.Properties
			err = ft.Geometry.UnmarshalJSONCoords(inft.Geometry.Type, inft.Geometry.Coordinates)
			if err != nil {
				return err
			}
			*fl = append(*fl, ft)
		} else {
			// Because the lib doesn't have native Multi* types, we split those into single geometries.
			var singles []json.RawMessage
			json.Unmarshal(inft.Geometry.Coordinates, &singles)
			for _, single := range singles {
				var ft spatial.Feature
				ft.Props = inft.Properties
				err = ft.Geometry.UnmarshalJSONCoords(inft.Geometry.Type[5:], single)
				if err != nil {
					return err
				}
				*fl = append(*fl, ft)
			}
		}
	}
	return nil
}

func singularType(typ string) bool {
	return !strings.HasPrefix(strings.ToLower(typ), "multi")
}
