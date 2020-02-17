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
	if len(fc.SRID) != 0 {
		geojsonFC := featureColl{}
		geojsonFC.Type = "FeatureCollection"
		geojsonFC.Features = fc.Features
		geojsonFC.CRS.Properties.Name, _ = sridOGC[fc.SRID]
		geojsonFC.CRS.Type = "name"
		return json.NewEncoder(w).Encode(&geojsonFC)
	}
	geojsonFC := featureCollNoCRS{}
	geojsonFC.Type = "FeatureCollection"
	geojsonFC.Features = fc.Features
	return json.NewEncoder(w).Encode(&geojsonFC)
}

func (c *Codec) Extensions() []string {
	return []string{"geojson", "json"}
}

type featureCollNoCRS struct {
	Type     string   `json:"type"`
	Features FeatList `json:"features"`
}

type featureColl struct {
	featureCollNoCRS
	CRS struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
}

type FeatureProto struct {
	Geometry struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	ID         string
	Properties map[string]interface{} `json:"properties"`
}

type FeatList []spatial.Feature

func (fl *FeatList) UnmarshalJSON(buf []byte) error {
	var fts []FeatureProto
	err := json.Unmarshal(buf, &fts)
	if err != nil {
		return err
	}

	*fl = make([]spatial.Feature, 0, len(fts))
	for _, inft := range fts {
		if inft.ID != "" {
			inft.Properties["id"] = inft.ID
		}

		err = fl.UnmarshalJSONCoords(inft)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fl *FeatList) UnmarshalJSONCoords(fp FeatureProto) error {
	var err error
	if singularType(fp.Geometry.Type) {
		var ft spatial.Feature
		ft.Props = fp.Properties
		err = ft.Geometry.UnmarshalJSONCoords(fp.Geometry.Type, fp.Geometry.Coordinates)
		if err == spatial.ErrorEmptyGeomType {
			// TODO: Shall we warn here somehow?
			return nil
		}
		if err != nil {
			return err
		}
		*fl = append(*fl, ft)
	} else {
		// Because the lib doesn't have native Multi* types, we split those into single geometries.
		var singles []json.RawMessage
		json.Unmarshal(fp.Geometry.Coordinates, &singles)
		for _, single := range singles {
			var ft spatial.Feature
			ft.Props = fp.Properties
			err = ft.Geometry.UnmarshalJSONCoords(fp.Geometry.Type[5:], single)
			if err != nil {
				return err
			}
			*fl = append(*fl, ft)
		}
	}
	return nil
}

func singularType(typ string) bool {
	return !strings.HasPrefix(strings.ToLower(typ), "multi")
}
