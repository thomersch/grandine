package main

import (
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

type FeatureCache interface {
	AddFeature(spatial.Feature)
	GetFeatures(tile.ID) []spatial.Feature

	BBox() spatial.BBox
	Count() int
}

type FeatureTable struct {
	Zoomlevels []int

	count int
	table map[int][][][]spatial.Feature
	bbox  *spatial.BBox
}

func NewFeatureTable(zoomlevels []int) *FeatureTable {
	ftab := FeatureTable{Zoomlevels: zoomlevels}

	ftab.table = map[int][][][]spatial.Feature{}
	for _, zl := range ftab.Zoomlevels {
		l := pow(2, zl)
		ftab.table[zl] = make([][][]spatial.Feature, l)
		for x := range ftab.table[zl] {
			ftab.table[zl][x] = make([][]spatial.Feature, l)
		}
	}
	return &ftab
}

func (ftab *FeatureTable) AddFeature(ft spatial.Feature) {
	for _, zl := range ftab.Zoomlevels {
		if !renderable(ft.Props, zl) {
			continue
		}
		for _, tid := range tile.Coverage(ft.Geometry.BBox(), zl) {
			ftab.table[zl][tid.X][tid.Y] = append(ftab.table[zl][tid.X][tid.Y], ft)
		}
	}
	if ftab.bbox == nil {
		var bb = ft.Geometry.BBox()
		ftab.bbox = &bb
	} else {
		ftab.bbox.ExtendWith(ft.Geometry.BBox())
	}

	ftab.count++
}

func (ftab *FeatureTable) GetFeatures(tid tile.ID) []spatial.Feature {
	return ftab.table[tid.Z][tid.X][tid.Y]
}

func (ftab *FeatureTable) Count() int {
	return ftab.count
}

func (ftab *FeatureTable) BBox() spatial.BBox {
	if ftab.bbox == nil {
		return spatial.BBox{}
	}
	return *ftab.bbox
}
