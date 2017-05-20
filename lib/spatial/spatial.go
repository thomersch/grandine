package spatial

import "encoding/json"

type PropertyRetriever interface {
	Properties() map[string]interface{}
}

// Feature is a data structure which holds geometry and tags/properties of a geographical feature.
type Feature struct {
	Type     string                 `json:"type"` //TODO: necessary?
	Props    map[string]interface{} `json:"properties"`
	Geometry Geom                   `json:"geometry"`
}

func (f *Feature) Properties() map[string]interface{} {
	return f.Props
}

func (f *Feature) MarshalWKB() ([]byte, error) {
	return f.Geometry.MarshalWKB()
}

type FeatureCollection struct {
	Features []Feature `json:"features"`
}

func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	wfc := struct {
		Type     string    `json:"type"`
		Features []Feature `json:"features"`
	}{
		Type:     "FeatureCollection",
		Features: fc.Features,
	}
	return json.Marshal(wfc)
}
