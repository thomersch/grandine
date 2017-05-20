package spatial

import "math"

func ringBBox(ring []Point) (nw, se Point) {
	nw[0] = ring[0][0]
	nw[1] = ring[0][1]
	se[0] = nw[0]
	se[1] = nw[1]
	for _, pt := range ring {
		nw[0] = math.Min(nw[0], pt[0])
		nw[1] = math.Min(nw[1], pt[1])
		se[0] = math.Max(se[0], pt[0])
		se[1] = math.Max(se[1], pt[1])
	}
	return
}
