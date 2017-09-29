package spatial

import (
	"encoding/json"
	"math"

	"github.com/dhconnelly/rtreego"
)

type PropertyRetriever interface {
	Properties() map[string]interface{}
}

type Filterable interface {
	Filter(BBox) []Feature
}

// Feature is a data structure which holds geometry and tags/properties of a geographical feature.
type Feature struct {
	Props    map[string]interface{} `json:"properties"`
	Geometry Geom
}

func NewFeature() Feature {
	return Feature{
		Props: map[string]interface{}{},
	}
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

func bboxToRect(bbox BBox) *rtreego.Rect {
	dist := []float64{bbox.NE.X() - bbox.SW.X(), bbox.NE.Y() - bbox.SW.Y()}
	// rtreego doesn't allow zero sized bboxes
	if dist[0] == 0 {
		dist[0] = math.SmallestNonzeroFloat64
	}
	if dist[1] == 0 {
		dist[1] = math.SmallestNonzeroFloat64
	}
	r, err := rtreego.NewRect(rtreego.Point{bbox.SW.X(), bbox.SW.Y()}, dist)
	if err != nil {
		panic(err)
	}
	return r
}

type FeatureCollection struct {
	Features []Feature `json:"features"`
	SRID     string    `json:"-"`
} // TODO: consider adding properties field

// Deprecated. Please use geojson.Codec.
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

func (fc *FeatureCollection) Filter(bbox BBox) []Feature {
	var filtered []Feature

	for _, feat := range fc.Features {
		if feat.Geometry.In(bbox) {
			filtered = append(filtered, feat)
		}
	}
	return filtered
}

type rtreeFeat Feature

func (ft rtreeFeat) Bounds() *rtreego.Rect {
	return bboxToRect(ft.Geometry.BBox())
}

// RTreeCollection is a FeatureCollection which is backed by a rtree.
type RTreeCollection struct {
	rt *rtreego.Rtree
}

func NewRTreeCollection(features ...Feature) *RTreeCollection {
	var fts []rtreego.Spatial
	for _, ft := range features {
		fts = append(fts, rtreeFeat(ft))
	}

	return &RTreeCollection{
		// TODO: find out optimal branching factor
		rt: rtreego.NewTree(2, 32, 64, fts...),
	}
}

func (rt *RTreeCollection) Add(feature Feature) {
	rt.rt.Insert(rtreeFeat(feature))
}

func (rt *RTreeCollection) Filter(bbox BBox) []Feature {
	var fts []Feature
	for _, ft := range rt.rt.SearchIntersect(bboxToRect(bbox)) {
		fts = append(fts, Feature(ft.(rtreeFeat)))
	}
	return fts
}
