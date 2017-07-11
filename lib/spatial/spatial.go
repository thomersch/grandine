package spatial

import (
	"encoding/json"

	"github.com/dhconnelly/rtreego"
)

type PropertyRetriever interface {
	Properties() map[string]interface{}
}

// Feature is a data structure which holds geometry and tags/properties of a geographical feature.
type Feature struct {
	Props    map[string]interface{} `json:"properties"`
	Geometry Geom
}

func (f *Feature) Properties() map[string]interface{} {
	return f.Props
}

func (f *Feature) MarshalWKB() ([]byte, error) {
	return f.Geometry.MarshalWKB()
}

func (f Feature) MarshalJSON() ([]byte, error) {
	tfc := struct {
		Type     string                 `json:"type"`
		Props    map[string]interface{} `json:"properties"`
		Geometry Geom                   `json:"geometry"`
	}{
		Type:     "Feature",
		Props:    f.Props,
		Geometry: f.Geometry,
	}
	return json.Marshal(tfc)
}

// Bounds is here to satisfy rtreego's interface. This is not guaranteed to be stable.
func (f Feature) Bounds() *rtreego.Rect {
	bbox := f.Geometry.BBox()
	r, err := rtreego.NewRect(rtreego.Point{bbox.SW.X(), bbox.SW.Y()}, []float64{bbox.NE.X() - bbox.SW.X(), bbox.NE.Y() - bbox.SW.Y()})
	if err != nil {
		panic(err)
	}
	return r
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

// RTreeCollection is a FeatureCollection which is backed by a rtree.
type RTreeCollection struct {
	rt *rtreego.Rtree
}

func NewRTreeCollection() RTreeCollection {
	return RTreeCollection{
		// TODO: find out optimal branching factor
		rt: rtreego.NewTree(2, 25, 50),
	}
}

func (rtc *RTreeCollection) UnmarshalJSON(buf []byte) error {
	// TODO: consider parsing progressively
	fcoll := FeatureCollection{}
	if err := json.Unmarshal(buf, &fcoll); err != nil {
		return err
	}
	for _, ft := range fcoll.Features {
		rtc.rt.Insert(ft)
	}
	return nil
}

func Filter(fcs []Feature, bbox BBox) []Feature {
	var filtered []Feature

	for _, feat := range fcs {
		if feat.Geometry.In(bbox) {
			filtered = append(filtered, feat)
		}
	}
	return filtered
}
