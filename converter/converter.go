package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/thomersch/grandine/fileformat"

	"github.com/golang/protobuf/proto"
	geom "github.com/twpayne/go-geom"
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
	Type       string
	Properties map[string]string
	Geometry   struct {
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

type Feature interface {
	WKB() []byte
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
			writeBlock(out, []Feature{&feat})
		}
	} else {
		log.Println("Decoding .unnamed, writing geojson")
		readFileHeader(f)
		// readBlocks()
	}
}

func writeFileHeader(w io.Writer) {
	// Cookie
	w.Write([]byte(cookie))

	// Version
	binary.Write(w, binary.LittleEndian, uint32(version))
}

func writeBlock(w io.Writer, fs []Feature) {
	blockBody := &fileformat.Body{}
	for _, f := range fs {
		blockBody.Feature = append(blockBody.Feature, &fileformat.Feature{
			Geom: f.WKB(),
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
