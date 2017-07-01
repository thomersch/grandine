package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thomersch/grandine/lib/cugdf"
	"github.com/thomersch/grandine/lib/mvt"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

var zoomlevels = []int{6, 7, 8, 9, 10, 11, 12, 13, 14}

func main() {
	source := flag.String("src", "geo.geojson", "file to read from, supported formats: geojson, cugdf")
	target := flag.String("target", "tiles", "path where the tiles will be written")
	flag.Parse()

	f, err := os.Open(*source)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = os.MkdirAll(*target, 0777)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("parsing input...")
	fc := spatial.FeatureCollection{}

	if strings.HasSuffix(strings.ToLower(*source), "geojson") {
		if err := json.NewDecoder(f).Decode(&fc); err != nil {
			log.Fatal(err)
		}
	} else {
		fc.Features, err = cugdf.Unmarshal(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("read %d features", len(fc.Features))

	var bboxPts []spatial.Point
	for _, feat := range fc.Features {
		bb := feat.Geometry.BBox()
		bboxPts = append(bboxPts, bb.SW, bb.NE)
	}

	bbox := spatial.Line(bboxPts).BBox()
	log.Println("filtering features...")

	features := spatial.Filter(fc.Features, bbox)
	if len(features) == 0 {
		log.Println("no features to be processed, exiting.")
		os.Exit(2)
	}
	log.Printf("%d features to be processed", len(features))

	var tc []tile.ID
	for _, zoomlevel := range zoomlevels {
		tc = append(tc, tile.Coverage(bbox, zoomlevel)...)
	}
	log.Printf("attempting to generate %d tiles", len(tc))

	for _, tID := range tc {
		log.Printf("Generating %v", tID)
		var tileFeatures []spatial.Feature
		tileClipBBox := tID.BBox()
		for _, feat := range spatial.Filter(features, tileClipBBox) {
			for _, geom := range feat.Geometry.ClipToBBox(tileClipBBox) {
				feat.Geometry = geom
				tileFeatures = append(tileFeatures, feat)
			}
		}
		if len(tileFeatures) == 0 {
			continue
		}
		layer := map[string][]spatial.Feature{
			"default": tileFeatures,
		}
		buf, err := mvt.EncodeTile(layer, tID)
		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(filepath.Join(*target, strconv.Itoa(tID.Z), strconv.Itoa(tID.X)), 0777)
		if err != nil {
			log.Fatal(err)
		}
		tf, err := os.Create(filepath.Join(*target, strconv.Itoa(tID.Z), strconv.Itoa(tID.X), strconv.Itoa(tID.Y)+".mvt"))
		if err != nil {
			log.Fatal(err)
		}
		_, err = tf.Write(buf)
		if err != nil {
			log.Fatal(err)
		}
		tf.Close()
	}
}
