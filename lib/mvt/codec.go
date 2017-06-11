package mvt

import (
	"errors"

	"github.com/golang/protobuf/proto"
	vt "github.com/thomersch/grandine/lib/mvt/vector_tile"
	"github.com/thomersch/grandine/lib/spatial"
)

type cmd uint32

const (
	cmdMoveTo    cmd = 1
	cmdLineTo    cmd = 2
	cmdClosePath cmd = 7

	extent = 4096
)

var (
	vtPoint        = vt.Tile_POINT
	vtLine         = vt.Tile_LINESTRING
	vtPoly         = vt.Tile_POLYGON
	vtLayerVersion = uint32(2)
)

func encodeCommandInt(c cmd, count uint32) uint32 {
	return (uint32(c) & 0x7) | (count << 3)
}

// func decodeCommandInt(cmdInt uint32) (cmd, uint32) {
// 	return cmd(cmdInt & 0x7), cmdInt >> 3
// }

func encodeZigZag(i int) uint32 {
	return uint32((i << 1) ^ (i >> 31))
}

func EncodeTile(features map[string][]spatial.Feature, tid TileID) ([]byte, error) {
	tile, err := assembleTile(features, tid)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&tile)
}

func assembleTile(features map[string][]spatial.Feature, tid TileID) (vt.Tile, error) {
	var tile vt.Tile
	for layerName, layerFeats := range features {
		layer, err := assembleLayer(layerFeats, tid)
		if err != nil {
			return tile, err
		}
		layer.Name = &layerName
		layer.Version = &vtLayerVersion
		tile.Layers = append(tile.Layers, &layer)
	}
	return tile, nil
}

func assembleLayer(features []spatial.Feature, tile TileID) (vt.Tile_Layer, error) {
	var (
		tl  vt.Tile_Layer
		err error
		ext = uint32(4096)
	)

	// TODO: tags

	for _, feat := range features {
		var tileFeat vt.Tile_Feature

		tileFeat.Geometry, err = encodeGeometry([]spatial.Geom{feat.Geometry}, tile)
		if err != nil {
			return tl, err
		}
		switch feat.Geometry.Typ() {
		case spatial.GeomTypePoint:
			tileFeat.Type = &vtPoint
		case spatial.GeomTypeLineString:
			tileFeat.Type = &vtLine
		case spatial.GeomTypePolygon:
			tileFeat.Type = &vtPoly
		default:
			return tl, errors.New("unknown geometry type")
		}

		tl.Features = append(tl.Features, &tileFeat)
		tl.Extent = &ext //TODO: configurable?
	}
	return tl, nil
}

// encodes one or more geometries of the same type into one (multi-)geometry
func encodeGeometry(geoms []spatial.Geom, tile TileID) (commands []uint32, err error) {
	var (
		cur    [2]uint32
		dx, dy uint32
		// the following four lines might be optimized
		nw, se           = tile.BBox()
		bbox             = bbox{nw, se}
		xScale, yScale   = tileScalingFactor(bbox, extent)
		xOffset, yOffset = tileOffset(bbox)
	)
	var typ spatial.GeomType
	for _, geom := range geoms {
		if typ != 0 && typ != geom.Typ() {
			return nil, errors.New("encodeGeometry only accepts uniform geoms")
		}
		switch geom.Typ() {
		case spatial.GeomTypePoint:
			pt, _ := geom.Point()
			tX, tY := tileCoord(pt, extent, xScale, yScale, xOffset, yOffset)
			dx = encodeZigZag(tX - int(cur[0]))
			dy = encodeZigZag(tY - int(cur[1]))
			commands = append(commands, encodeCommandInt(cmdMoveTo, 1), dx, dy)
		}
	}
	return commands, nil
}
