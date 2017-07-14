package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"

	"github.com/thomersch/grandine/lib/cugdf"
	"github.com/thomersch/grandine/lib/mvt"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

type zmLvl []int

func (zm *zmLvl) String() string {
	return fmt.Sprintf("%d", *zm)
}

func (zm *zmLvl) Set(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s (only integer values are allowed)", value)
	}
	*zm = append(*zm, v)
	return nil
}

var zoomlevels zmLvl

func main() {
	source := flag.String("src", "geo.geojson", "file to read from, supported formats: geojson, cugdf")
	target := flag.String("target", "tiles", "path where the tiles will be written")
	defaultLayer := flag.Bool("default-layer", true, "if no layer name is specified in the feature, whether it will be put into a default layer")
	workersNumber := flag.Int("workers", runtime.GOMAXPROCS(0), "number of workers")
	cpuProfile := flag.String("cpuprof", "", "writes CPU profiling data into a file")
	flag.Var(&zoomlevels, "zoom", "one or more zoom level of which the tiles will be rendered")
	flag.Parse()

	if len(zoomlevels) == 0 {
		log.Fatal("no zoom levels specified")
	}

	if len(*cpuProfile) != 0 {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	if len(fc.Features) == 0 {
		log.Fatal("no features in input file")
	}

	log.Printf("read %d features", len(fc.Features))

	var (
		bboxPts []spatial.Point
		rt      = spatial.NewRTreeCollection()
	)
	for _, feat := range fc.Features {
		bb := feat.Geometry.BBox()
		bboxPts = append(bboxPts, bb.SW, bb.NE)
		rt.Add(feat)
	}

	bbox := spatial.Line(bboxPts).BBox()
	log.Println("filtering features...")

	// features := spatial.Filter(fc.Features, bbox)
	// if len(features) == 0 {
	// 	log.Println("no features to be processed, exiting.")
	// 	os.Exit(2)
	// }
	// log.Printf("%d features to be processed", len(features))

	var tc []tile.ID
	for _, zoomlevel := range zoomlevels {
		tc = append(tc, tile.Coverage(bbox, zoomlevel)...)
	}
	log.Printf("attempting to generate %d tiles", len(tc))

	dtw := diskTileWriter{basedir: *target}
	dlm := defaultLayerMapper{defaultLayer: *defaultLayer}

	var (
		wg sync.WaitGroup
		ws = workerSlices(tc, *workersNumber)
	)
	for wrk := 0; wrk < len(ws); wrk++ {
		wg.Add(1)
		go func(i int) {
			generateTiles(ws[i], rt, &dtw, &dlm)
			wg.Done()
		}(wrk)
	}
	wg.Wait()
}

func workerSlices(tiles []tile.ID, wrkNum int) [][]tile.ID {
	var r [][]tile.ID
	if len(tiles) <= wrkNum {
		for t := 0; t < len(tiles); t++ {
			r = append(r, []tile.ID{tiles[t]})
		}
		return r
	}
	for wrkr := 0; wrkr < wrkNum; wrkr++ {
		start := (len(tiles) / wrkNum) * wrkr
		end := (len(tiles) / wrkNum) * (wrkr + 1)
		if wrkr == wrkNum-1 {
			end = len(tiles)
		}
		r = append(r, tiles[start:end])
	}
	return r
}

type diskTileWriter struct {
	basedir string
}

func (tw *diskTileWriter) WriteTile(tID tile.ID, buf []byte) error {
	err := os.MkdirAll(filepath.Join(tw.basedir, strconv.Itoa(tID.Z), strconv.Itoa(tID.X)), 0777)
	if err != nil {
		return err
	}
	tf, err := os.Create(filepath.Join(tw.basedir, strconv.Itoa(tID.Z), strconv.Itoa(tID.X), strconv.Itoa(tID.Y)+".mvt"))
	if err != nil {
		return err
	}
	defer tf.Close()
	_, err = tf.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

type defaultLayerMapper struct {
	defaultLayer bool
}

func (dlm *defaultLayerMapper) LayerName(props map[string]interface{}) string {
	if layerName, ok := props["@layer"]; ok {
		return layerName.(string)
	}
	if dlm.defaultLayer {
		return "default"
	}
	return ""
}

type layerMapper interface {
	LayerName(map[string]interface{}) string
}

type tileWriter interface {
	WriteTile(tile.ID, []byte) error
}

func generateTiles(tIDs []tile.ID, features spatial.RTreeCollection, tw tileWriter, lm layerMapper) {
	for _, tID := range tIDs {
		log.Printf("Generating %s", tID)
		var (
			layers = map[string][]spatial.Feature{}
			ln     string
		)
		tileClipBBox := tID.BBox()
		// for _, feat := range spatial.Filter(features, tileClipBBox) {
		for _, feat := range features.Find(tileClipBBox) {
			for _, geom := range feat.Geometry.ClipToBBox(tileClipBBox) {
				feat.Geometry = geom
				ln = lm.LayerName(feat.Props)
				if len(ln) != 0 {
					if _, ok := layers[ln]; !ok {
						layers[ln] = []spatial.Feature{feat}
					} else {
						layers[ln] = append(layers[ln], feat)
					}
				}
			}
		}
		if !anyFeatures(layers) {
			continue
		}
		buf, err := mvt.EncodeTile(layers, tID)
		if err != nil {
			log.Fatal(err)
		}

		err = tw.WriteTile(tID, buf)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func anyFeatures(layers map[string][]spatial.Feature) bool {
	for _, ly := range layers {
		if len(ly) > 0 {
			return true
		}
	}
	return false
}
