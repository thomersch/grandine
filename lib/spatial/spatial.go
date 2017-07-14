package spatial

import (
	"encoding/json"
	"math"

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

type rtreeFeat Feature

func (ft rtreeFeat) Bounds() *rtreego.Rect {
	return bboxToRect(ft.Geometry.BBox())
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

func (rt *RTreeCollection) Add(feature Feature) {
	rt.rt.Insert(rtreeFeat(feature))
}

func (rt *RTreeCollection) Find(bbox BBox) []Feature {
	var fts []Feature
	for _, ft := range rt.rt.SearchIntersect(bboxToRect(bbox)) {
		fts = append(fts, Feature(ft.(rtreeFeat)))
	}
	return fts
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
