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

	extent      = 4096
	simplFactor = 2
)

var (
	vtPoint        = vt.Tile_POINT
	vtLine         = vt.Tile_LINESTRING
	vtPoly         = vt.Tile_POLYGON
	vtLayerVersion = uint32(2)

	skipAtKeys = true // if enabled keys that start with "@" will be ignored

	errNoGeom = errors.New("no valid geometries")
)

type Codec struct{}

func (c *Codec) EncodeTile(features map[string][]spatial.Feature, tid tile.ID) ([]byte, error) {
	return EncodeTile(features, tid)
}

func (c *Codec) Extension() string {
	return "mvt"
}

func encodeCommandInt(c cmd, count uint32) uint32 {
	return (uint32(c) & 0x7) | (count << 3)
}

func decodeCommandInt(cmdInt uint32) (cmd, uint32) {
	return cmd(cmdInt & 0x7), cmdInt >> 3
}

func encodeZigZag(i int) uint32 {
	return uint32((i << 1) ^ (i >> 31))
}

func EncodeTile(features map[string][]spatial.Feature, tid tile.ID) ([]byte, error) {
	vtile, err := assembleTile(features, tid)
	if err != nil {
		return nil, err
	}
	if len(vtile.Layers) == 0 {
		log.Println("no layers")
		return nil, nil
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
		if len(layer.Features) == 0 {
			log.Println("no features")
			continue
		}
		var ln = layerName
		layer.Name = &ln // &layerName can't be used directly, because pointers are reused in for range loops
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
		ext  = uint32(extent)
		keys = tagElems{}
		vals = tagElems{}
	)

	for _, feat := range features {
		var tileFeat vt.Tile_Feature

		for k, v := range feat.Properties() {
			if skipAtKeys && k[0] == '@' {
				continue
			}
			kpos := keys.Index(k)
			vpos := vals.Index(v)
			tileFeat.Tags = append(tileFeat.Tags, uint32(kpos), uint32(vpos))
		}

		tileFeat.Geometry, err = encodeGeometry([]spatial.Geom{feat.Geometry}, tid)
		if len(tileFeat.Geometry) == 0 || err == errNoGeom {
			continue
		}
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
		tp   tileParams
		bbox = tid.BBox()
	)
	tp.xScale, tp.yScale = tileScalingFactor(bbox, extent)
	tp.xOffset, tp.yOffset = tileOffset(bbox)
	tp.extent = extent
	tileBox := spatial.BBox{spatial.Point{0, 0}, spatial.Point{float64(tp.extent) * 1.05, float64(tp.extent) * 1.05}}

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
			tX, tY := tileCoord(pt, tp)
			if tX > extent || tY > extent {
				log.Printf("%v is outside of tile", pt)
			}
			dx = tX - int(cur[0])
			dy = tY - int(cur[1])
			// TODO: support multipoint
			commands = append(commands, encodeCommandInt(cmdMoveTo, 1), encodeZigZag(dx), encodeZigZag(dy))
		case spatial.GeomTypeLineString:
			ln := simplifyAndClipInTile(geom.MustLineString(), tileBox, tp)
			if ln == nil {
				continue
			}
			commands = append(commands, encodeLine(ln, &cur)...)
		case spatial.GeomTypePolygon:
			poly, _ := geom.Polygon()
			log.Println(poly)
			for i, ring := range poly {
				ring = lineToTileCoords(ring, tp)
				if ring[0] == ring[len(ring)-1] {
					ring = ring[:len(ring)-1]
				}
				// ring = ring.Simplify(simplFactor)
				// if len(ring) == 0 {
				// 	continue
				// }
				poly[i] = ring
				// gms := ring.ClipToBBox(tileBox)
				// if len(gms) == 0 || len(gms[0].MustLineString()) == 0 {
				// 	continue
				// }
				// poly[i], _ = gms[0].LineString()
			}

			if !poly.Valid() {
				log.Println("invalid")
				continue
			}
			poly = poly.Simplify(simplFactor)

			for _, poly := range poly.ClipToBBox(tileBox) {
				for _, ring := range poly.MustPolygon() {
					l := encodeLine(ring, &cur)
					if l == nil {
						return nil, errNoGeom
					}
					commands = append(commands, l...)
					commands = append(commands, encodeCommandInt(cmdClosePath, 1))
				}
			}
		}
	}
	return commands, nil
}

func simplifyAndClipInTile(ln spatial.Line, bbox spatial.BBox, tp tileParams) spatial.Line {
	ln = lineToTileCoords(ln, tp)
	ln = ln.Simplify(simplFactor)
	if len(ln) == 0 {
		return nil
	}
	gm := ln.ClipToBBox(bbox)
	if len(gm) == 0 {
		return nil
	}
	ln = gm[0].MustLineString()
	if len(ln) == 0 {
		return nil
	}
	return ln
}

// encodeLine takes a line in tile coordinates
func encodeLine(ln spatial.Line, cur *[2]int) []uint32 {
	if len(ln) == 0 {
		return nil
	}
	var (
		commands = make([]uint32, len(ln)*2+2) // len=number of coordinates + initial move to + size
		dx, dy   int
	)
	commands[0] = encodeCommandInt(cmdMoveTo, 1)
	commands[3] = encodeCommandInt(cmdLineTo, uint32(len(commands)-4)/2)
	for i, tc := range ln {
		dx = int(tc.X) - cur[0]
		dy = int(tc.Y) - cur[1]
		cur[0] = int(tc.X)
		cur[1] = int(tc.Y)
		if i == 0 {
			commands[1] = encodeZigZag(int(dx))
			commands[2] = encodeZigZag(int(dy))
		} else {
			commands[i+i+2] = encodeZigZag(int(dx))
			commands[i+i+3] = encodeZigZag(int(dy))
		}
	}
	return commands
}

func lineToTileCoords(ln spatial.Line, tp tileParams) spatial.Line {
	if len(ln) == 0 {
		return nil
	}
	var (
		tlCrds = make(spatial.Line, 0, len(ln))
		tx, ty int
	)
	for _, pt := range ln {
		tx, ty = tileCoord(pt, tp)
		tlCrds = append(tlCrds, spatial.Point{float64(tx), float64(ty)})
	}
	return tlCrds
}
