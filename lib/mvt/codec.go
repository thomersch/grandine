package mvt

import (
	"errors"
	"fmt"
	"log"

	vt "github.com/thomersch/grandine/lib/mvt/vector_tile"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"

	"github.com/golang/protobuf/proto"
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

func EncodeTile(features map[string][]spatial.Feature, tid tile.ID) ([]byte, error) {
	vtile, err := assembleTile(features, tid)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&vtile)
}

func assembleTile(features map[string][]spatial.Feature, tid tile.ID) (vt.Tile, error) {
	var vtile vt.Tile
	for layerName, layerFeats := range features {
		layer, err := assembleLayer(layerFeats, tid)
		if err != nil {
			return vtile, err
		}
		layer.Name = &layerName
		layer.Version = &vtLayerVersion
		vtile.Layers = append(vtile.Layers, &layer)
	}
	return vtile, nil
}

// tagElems is an intermediate data structure for serializing keys or values into flat
// index referenced lists as used by MVT
type tagElems map[interface{}]int

func (te tagElems) Index(v interface{}) int {
	if pos, ok := te[v]; ok {
		return pos
	}
	pos := len(te)
	te[v] = pos
	return pos
}

func (te tagElems) Strings() []string {
	var l = make([]string, len(te))
	for elem, pos := range te {
		l[pos] = elem.(string)
	}
	return l
}

func (te tagElems) Values() []*vt.Tile_Value {
	var l = make([]*vt.Tile_Value, len(te))
	for val, pos := range te {
		var tv vt.Tile_Value
		switch v := val.(type) {
		case string:
			tv.StringValue = &v
		case float32:
			tv.FloatValue = &v
		case float64:
			tv.DoubleValue = &v
		case int:
			i := int64(v)
			tv.SintValue = &i
		case int64:
			tv.SintValue = &v
		case uint:
			i := uint64(v)
			tv.UintValue = &i
		case uint64:
			tv.UintValue = &v
		case bool:
			tv.BoolValue = &v
		default:
			s := fmt.Sprintf("%s", v)
			tv.StringValue = &s
		}
		l[pos] = &tv
	}
	return l
}

func assembleLayer(features []spatial.Feature, tid tile.ID) (vt.Tile_Layer, error) {
	var (
		tl   vt.Tile_Layer
		err  error
		ext  = uint32(4096)
		keys = tagElems{}
		vals = tagElems{}
	)

	for _, feat := range features {
		var tileFeat vt.Tile_Feature

		for k, v := range feat.Properties() {
			kpos := keys.Index(k)
			vpos := vals.Index(v)
			tileFeat.Tags = append(tileFeat.Tags, uint32(kpos), uint32(vpos))
		}

		tileFeat.Geometry, err = encodeGeometry([]spatial.Geom{feat.Geometry}, tid)
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

	tl.Keys = keys.Strings()
	tl.Values = vals.Values()
	return tl, nil
}

// encodes one or more geometries of the same type into one (multi-)geometry
func encodeGeometry(geoms []spatial.Geom, tid tile.ID) (commands []uint32, err error) {
	var (
		cur    [2]int
		dx, dy int
		// the following four lines might be optimized
		sw, ne           = tid.BBox()
		bbox             = bbox{sw, ne}
		xScale, yScale   = tileScalingFactor(bbox, extent)
		xOffset, yOffset = tileOffset(bbox)
	)
	var typ spatial.GeomType
	for _, geom := range geoms {
		if typ != 0 && typ != geom.Typ() {
			return nil, errors.New("encodeGeometry only accepts uniform geoms")
		}

		cur[0] = 0
		cur[1] = 0
		switch geom.Typ() {
		case spatial.GeomTypePoint:
			pt, _ := geom.Point()
			tX, tY := tileCoord(pt, extent, xScale, yScale, xOffset, yOffset)
			if tX > extent || tY > extent {
				log.Printf("%v is outside of tile", pt)
			}
			dx = tX - int(cur[0])
			dy = tY - int(cur[1])
			// TODO: support multipoint
			commands = append(commands, encodeCommandInt(cmdMoveTo, 1), encodeZigZag(dx), encodeZigZag(dy))
		case spatial.GeomTypeLineString:
			ln, _ := geom.LineString()
			commands = append(commands, encodeLine(ln, cur, extent, xScale, yScale, xOffset, yOffset)...)
		case spatial.GeomTypePolygon:
			// TODO: inner ring handling
			poly, _ := geom.Polygon()
			commands = append(commands, encodeLine(poly[0], cur, extent, xScale, yScale, xOffset, yOffset)...)
			commands = append(commands, encodeCommandInt(cmdClosePath, 1))
		}
	}
	return commands, nil
}

func encodeLine(ln spatial.Line, cur [2]int, extent int, xScale, yScale, xOffset, yOffset float64) (commands []uint32) {
	tX, tY := tileCoord(ln[0], extent, xScale, yScale, xOffset, yOffset)
	dx := tX - int(cur[0])
	dy := tY - int(cur[1])
	cur[0] = cur[0] + dx
	cur[1] = cur[1] + dy

	commands = append(commands, encodeCommandInt(cmdMoveTo, 1), encodeZigZag(dx), encodeZigZag(dy),
		encodeCommandInt(cmdLineTo, uint32(len(ln)-1)))

	for _, pt := range ln[1:] {
		tX, tY = tileCoord(pt, extent, xScale, yScale, xOffset, yOffset)
		dx = tX - int(cur[0])
		dy = tY - int(cur[1])
		commands = append(commands, encodeZigZag(dx), encodeZigZag(dy))
		cur[0] = cur[0] + dx
		cur[1] = cur[1] + dy
	}
	return commands
}
