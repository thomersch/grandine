package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/thomersch/grandine/converter/fileformat"
	"github.com/thomersch/grandine/lib/cugdf"
	"github.com/thomersch/grandine/lib/spatial"

	"github.com/golang/protobuf/proto"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

const (
	cookie  = "ABCD"
	version = 0
)

type GeoJSON struct {
	Type     string
	Features []geoJSONFeature
}

type geoJSONFeature struct {
	Type     string
	Props    map[string]interface{} `json:"properties"`
	Geometry struct {
		Type        string
		Coordinates json.RawMessage
	}
}

func (g *geoJSONFeature) WKB() []byte {
	var (
		err       error
		elem      geom.T
		unmarshal = func(c interface{}) {
			err = json.Unmarshal(g.Geometry.Coordinates, c)
		}
	)

	switch g.Geometry.Type {
	case "Point":
		var coords geom.Coord
		unmarshal(&coords)
		elem = geom.NewPoint(geom.XY).MustSetCoords(coords)
	case "LineString":
		var coords []geom.Coord
		unmarshal(&coords)
		elem = geom.NewLineString(geom.XY).MustSetCoords(coords)
	case "Polygon":
		var coords [][]geom.Coord
		unmarshal(&coords)
		elem = geom.NewPolygon(geom.XY).MustSetCoords(coords)
	case "MultiPolygon":
		var coords [][][]geom.Coord
		unmarshal(&coords)
		elem = geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)
	default:
		log.Fatalf("type %v not yet implemented", g.Geometry.Type)
	}
	if err != nil {
		log.Fatal(err)
	}
	buf, err := wkb.Marshal(elem, binary.LittleEndian)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

func (g *geoJSONFeature) Properties() map[string]interface{} {
	return g.Props
}

type GeoJSONable interface {
	Geom() geom.T
	Properties() map[string]interface{}
}

func main() {
	source := flag.String("src", "geo.geojson", "")
	dest := flag.String("dest", "geo.unnamed", "")
	flag.Parse()

	f, err := os.Open(*source)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	out, err := os.Create(*dest)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if strings.HasSuffix(*source, ".geojson") || strings.HasSuffix(*source, ".json") {
		log.Println("Decoding geojson, encoding .unnamed")
		err := cugdf.WriteFileHeader(out)
		if err != nil {
			log.Fatal(err)
		}

		fc := spatial.FeatureCollection{}
		if err := json.NewDecoder(f).Decode(&fc); err != nil {
			log.Fatal(err)
		}

		for _, featBlock := range geomBlocks(100, fc.Features) {
			cugdf.WriteBlock(out, featBlock)
		}
	} else {
		log.Println("Decoding .unnamed, writing geojson")
		readFileHeader(f)
		var gjfc geojson.FeatureCollection
		for _, ft := range readBlocks(f) {
			gjfc.Features = append(gjfc.Features, &geojson.Feature{
				Geometry:   ft.Geom(),
				Properties: ft.Properties(),
			})
		}
		buf, err := gjfc.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}
		out.Write(buf)
	}
}

// geomBlocks slices a slice of geometries into slices with a max size
func geomBlocks(size int, src []spatial.Feature) [][]spatial.Feature {
	if len(src) <= size {
		return [][]spatial.Feature{src}
	}

	var (
		i   int
		res [][]spatial.Feature
		end int
	)
	for end < len(src) {
		end = (i + 1) * size
		if end > len(src) {
			end = len(src)
		}
		res = append(res, src[i*size:end])
		i++
	}
	return res
}

func readFileHeader(r io.Reader) {
	var (
		ck   = make([]byte, 4)
		vers uint32
	)
	if _, err := r.Read(ck); err != nil {
		log.Fatal(err)
	}
	if string(ck) != cookie {
		log.Fatal("invalid cookie")
	}

	if err := binary.Read(r, binary.LittleEndian, &vers); err != nil {
		log.Fatal(err)
	}
	if vers > version {
		log.Fatal("invalid file version")
	}
}

type geomWithProps struct {
	geom  geom.T
	props map[string]interface{}
}

func (g *geomWithProps) Geom() geom.T {
	return g.geom
}

func (g *geomWithProps) Properties() map[string]interface{} {
	return g.props
}

func readBlocks(r io.Reader) []GeoJSONable {
	var fs []GeoJSONable
	for {
		var (
			blockLength uint32
			flags       uint16
			compression uint8
			messageType uint8
		)
		if err := binary.Read(r, binary.LittleEndian, &blockLength); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if err := binary.Read(r, binary.LittleEndian, &flags); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(r, binary.LittleEndian, &compression); err != nil {
			if compression != 0 {
				log.Fatal("compression is not supported")
			}
		}
		if err := binary.Read(r, binary.LittleEndian, &messageType); err != nil {
			if messageType != 0 {
				log.Fatal("message type is not supported")
			}
		}

		var (
			buf       = make([]byte, blockLength)
			blockBody fileformat.Body
		)
		if _, err := r.Read(buf); err != nil {
			log.Fatal(err)
		}
		if err := proto.Unmarshal(buf, &blockBody); err != nil {
			log.Fatal(err)
		}
		for _, f := range blockBody.GetFeature() {
			geom, err := wkb.Unmarshal(f.Geom)
			if err != nil {
				log.Fatal(err)
			}

			gwp := geomWithProps{geom: geom, props: make(map[string]interface{})}
			for _, tag := range f.Tags {
				k, v, err := fileformat.KeyValue(tag)
				if err != nil {
					// TODO
					log.Fatal(err)
				}
				gwp.props[k] = v
			}
			fs = append(fs, &gwp)
		}
	}
	return fs
}
