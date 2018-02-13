package tile

import (
	"bytes"

	"github.com/thomersch/grandine/lib/geojson"
	"github.com/thomersch/grandine/lib/spatial"
)

type Codec interface {
	EncodeTile(features map[string][]spatial.Feature, tid ID) ([]byte, error)
	Extension() string
}

// GeoJSONCodec for debugging purposes. Probably not suitable for any real-world use.
type GeoJSONCodec struct{}

func (g *GeoJSONCodec) EncodeTile(features map[string][]spatial.Feature, tid ID) ([]byte, error) {
	var buf bytes.Buffer
	var fts []spatial.Feature
	for _, ly := range features {
		for _, ft := range ly {
			fts = append(fts, ft)
		}
	}
	gjc := &geojson.Codec{}
	err := gjc.Encode(&buf, &spatial.FeatureCollection{Features: fts})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GeoJSONCodec) Extension() string {
	return "geojson"
}
