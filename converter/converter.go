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
		Coordinates [][][2]float64
	}
}

func (g *geoJSONFeature) WKB() []byte {
	var crdss [][]geom.Coord

	for _, coordset := range g.Geometry.Coordinates {
		var crds []geom.Coord
		for _, crd := range coordset {
			crds = append(crds, geom.Coord{crd[0], crd[1]})
		}
		crdss = append(crdss, crds)
	}

	var elem geom.T
	switch g.Geometry.Type {
	case "LineString":
		elem = geom.NewLineString(geom.XY).MustSetCoords(crdss[0])
	case "Polygon":
		elem = geom.NewPolygon(geom.XY).MustSetCoords(crdss)
	default:
		log.Fatalf("type %v not yet implemented", g.Geometry.Type)
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

type WKBable interface {
	WKB() []byte
	Properties() map[string]interface{}
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

	if strings.HasSuffix(*source, ".geojson") {
		log.Println("Decoding geojson, encoding .unnamed")
		writeFileHeader(out)

		var gj GeoJSON
		if err := json.NewDecoder(f).Decode(&gj); err != nil {
			log.Fatal(err)
		}

		// for now write a block for every feature
		for _, feat := range gj.Features {
			writeBlock(out, []WKBable{&feat})
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

func writeFileHeader(w io.Writer) {
	// Cookie
	w.Write([]byte(cookie))

	// Version
	binary.Write(w, binary.LittleEndian, uint32(version))
}

func writeBlock(w io.Writer, fs []WKBable) {
	blockBody := &fileformat.Body{}
	for _, f := range fs {
		var (
			tags []*fileformat.Tag
		)
		for k, v := range f.Properties() {
			val, typ, err := fileformat.ValueType(v)
			if err != nil {
				// TODO: convert to error reporting
				log.Fatal(err)
			}
			tags = append(tags, &fileformat.Tag{
				Key:   k,
				Value: val,
				Type:  typ,
			})
		}

		blockBody.Feature = append(blockBody.Feature, &fileformat.Feature{
			Geom: f.WKB(),
			Tags: tags,
		})
	}
	bodyBuf, err := proto.Marshal(blockBody)
	if err != nil {
		log.Fatal(err)
	}

	// Body Length (fill later)
	binary.Write(w, binary.LittleEndian, uint32(len(bodyBuf)))
	// Flags
	binary.Write(w, binary.LittleEndian, uint16(0))
	// Compression
	binary.Write(w, binary.LittleEndian, uint8(0))
	// Message Type
	binary.Write(w, binary.LittleEndian, uint8(0))

	w.Write(bodyBuf)
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
