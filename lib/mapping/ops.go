package mapping

import "github.com/thomersch/grandine/lib/spatial"

func polyToLines(g spatial.Geom) []spatial.Geom {
	poly, err := g.Polygon()
	if err != nil {
		return nil
	}
	var lines = make([]spatial.Geom, 0, len(poly))
	for _, ring := range poly {
		lines = append(lines, spatial.MustNewGeom(ring))
	}
	return lines
}
