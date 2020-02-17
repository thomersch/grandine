// +build !golangclip

package spatial

import (
	"log"

	"github.com/pmezard/gogeos/geos"
)

func (p Polygon) clipToBBox(b BBox) []Geom {
	gpoly := p.geos()
	if gpoly == nil {
		return nil
	}

	var bboxLine = make([]geos.Coord, 0, 4)
	for _, pt := range NewLinesFromSegments(BBoxBorders(b.SW, b.NE))[0] {
		bboxLine = append(bboxLine, geos.NewCoord(pt.X, pt.Y))
	}

	bboxPoly := geos.Must(geos.NewPolygon(bboxLine))
	res, err := bboxPoly.Intersection(gpoly)
	if err != nil {
		// Sometimes there is a minor topology problem, a zero buffer helps often.
		gpolyBuffed, err := gpoly.Buffer(0)
		if err != nil {
			panic(err)
		}
		res, err = bboxPoly.Intersection(gpolyBuffed)
		if err != nil {
			panic(err)
		}
	}

	var resGeoms []Geom
	for _, poly := range geosToPolygons(res) {
		resGeoms = append(resGeoms, MustNewGeom(poly))
	}
	return resGeoms
}

func (p Polygon) geos() *geos.Geometry {
	var rings = make([][]geos.Coord, 0, len(p))
	for _, ring := range p {
		var rg = make([]geos.Coord, 0, len(ring))
		for _, pt := range ring {
			rg = append(rg, geos.NewCoord(pt.X, pt.Y))
		}
		rg = append(rg, rg[0])
		rings = append(rings, rg)
	}
	var gpoly *geos.Geometry
	if len(rings) == 0 {
		return nil
	}
	if len(rings) > 1 {
		return geos.Must(geos.NewPolygon(rings[0], rings[1:]...))
	}
	gpoly, err := geos.NewPolygon(rings[0])
	if err != nil {
		log.Printf("invalid polygon: %v", err)
		return nil
	}
	return gpoly
}

func geosToPolygons(g *geos.Geometry) []Polygon {
	ty, _ := g.Type()
	if ty == geos.POLYGON {
		return []Polygon{geosToPolygon(g)}
	}
	nmax, err := g.NGeometry()
	if err != nil {
		panic(err)
	}
	var polys = make([]Polygon, 0, nmax)
	for n := 0; n < nmax; n++ {
		polys = append(polys, geosToPolygon(geos.Must(g.Geometry(n))))
	}
	return polys
}

func geosToPolygon(g *geos.Geometry) Polygon {
	sh, err := g.Shell()
	if err != nil {
		return Polygon{}
	}
	crds, err := sh.Coords()
	if err != nil {
		panic(err)
	}
	if len(crds) == 0 { // we got an empty polygon
		return Polygon{}
	}
	var (
		p    = make(Polygon, 0, 8)
		ring = make([]Point, 0, len(crds))
	)
	for _, crd := range crds {
		ring = append(ring, Point{crd.X, crd.Y})
	}
	p = append(p, ring[:len(ring)-1])

	holes, _ := g.Holes()
	for _, hole := range holes {
		crds, err = hole.Coords()
		if err != nil {
			panic(err)
		}
		ring = make([]Point, 0, len(crds))
		for _, crd := range crds {
			ring = append(ring, Point{crd.X, crd.Y})
		}
		p = append(p, ring[:len(ring)-1])
	}
	return p
}
