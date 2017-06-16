package tile

import "github.com/thomersch/grandine/lib/spatial"

func Coverage(bb spatial.BBox, zoom int) []ID {
	// Tiles are counted from top-left to bottom-right
	tl := spatial.Point{bb.SW.X(), bb.NE.Y()}
	br := spatial.Point{bb.NE.X(), bb.SW.Y()}

	p1 := TileName(tl, zoom)
	p2 := TileName(br, zoom)

	var tiles []ID

	for x := p1.X; x <= p2.X; x++ {
		for y := p1.Y; y <= p2.Y; y++ {
			tiles = append(tiles, ID{X: x, Y: y, Z: zoom})
		}
	}
	return tiles
}
