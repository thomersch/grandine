// +build golangclip

package spatial

import polyclip "github.com/ctessum/polyclip-go"

func (p Polygon) clipToBBox(b BBox) []Geom {
	pcp := p.polyclipPolygon()
	bboxPoly := polyclip.Polygon{polyclip.Contour{
		polyclip.Point{b.SW.X(), b.SW.Y()},
		polyclip.Point{b.NE.X(), b.SW.Y()},
		polyclip.Point{b.NE.X(), b.NE.Y()},
		polyclip.Point{b.SW.X(), b.NE.Y()},
	}}
	return []Geom{
		MustNewGeom(
			polyClipToPolygon(pcp.Construct(polyclip.INTERSECTION, bboxPoly)),
		),
	}
}

func (p Polygon) polyclipPolygon() polyclip.Polygon {
	// TODO: polyclip has the same data structure, check if there is any possiblity for some speed-up hack
	var pcp polyclip.Polygon
	for _, ring := range p {
		var cnt polyclip.Contour
		for _, pt := range ring {
			cnt = append(cnt, polyclip.Point{pt.X(), pt.Y()})
		}
		pcp = append(pcp, cnt)
	}
	return pcp
}

func polyClipToPolygon(pcp polyclip.Polygon) Polygon {
	var p Polygon
	for _, countour := range pcp {
		var ring Line
		for _, pt := range countour {
			ring = append(ring, Point{pt.X, pt.Y})
		}
		p = append(p, ring)
	}
	return p
}
