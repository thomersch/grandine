package tile

import "github.com/thomersch/grandine/lib/spatial"

func Coverage(bbs []spatial.BBox, zoomlvl int) []ID {
	for _, bb := range bbs {
		tn1 := TileName(bb.SE, zoomlvl)
		tn2 := TileName(bb.NW, zoomlvl)
	}
}
