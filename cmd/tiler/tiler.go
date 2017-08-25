package main

import (
	"compress/gzip"
	"errors"
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
	"time"

	"github.com/thomersch/grandine/lib/mvt"
	"github.com/thomersch/grandine/lib/progressbar"
	"github.com/thomersch/grandine/lib/spaten"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

const indexThreshold = 100000000

type zmLvl []int

func (zm *zmLvl) String() string {
	return fmt.Sprintf("%d", *zm)
}

func (zm *zmLvl) Set(value string) error {
	for _, s := range strings.Split(value, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return fmt.Errorf("%s (only integer values are allowed)", value)
		}
		*zm = append(*zm, v)
	}
	return nil
}

type bbox spatial.BBox

func (b *bbox) String() string {
	return fmt.Sprintf("%v", *b)
}

func (b *bbox) Set(value string) error {
	var (
		fl    [4]float64
		parts = strings.Split(value, ",")
		err   error
	)
	if len(parts) != 4 {
		return errors.New("bbox takes exactly 4 parameters: SW Lon, SW Lat, NE Lon, NE Lat")
	}
	for i, s := range parts {
		fl[i], err = strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("could not parse bbox expression: %v", err)
		}
	}
	b.SW = spatial.Point{fl[0], fl[1]}
	b.NE = spatial.Point{fl[2], fl[3]}
	return nil
}

var (
	zoomlevels zmLvl
	quiet      *bool
)

func main() {
	source := flag.String("in", "geo.geojson", "file to read from, supported format: spaten")
	target := flag.String("out", "tiles", "path where the tiles will be written")
	defaultLayer := flag.Bool("default-layer", true, "if no layer name is specified in the feature, whether it will be put into a default layer")
	workersNumber := flag.Int("workers", runtime.GOMAXPROCS(0), "number of workers")
	cpuProfile := flag.String("cpuprof", "", "writes CPU profiling data into a file")
	compressTiles := flag.Bool("compress", false, "compress tiles with gzip")
	quiet = flag.Bool("q", false, "argument to use if program should be run in quiet mode with reduced logging")

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
	codec := spaten.Codec{}
	codec.Decode(f, &fc)

	if len(fc.Features) == 0 {
		log.Fatal("no features in input file")
	}

	log.Printf("read %d features", len(fc.Features))

	var bboxPts []spatial.Point
	for _, feat := range fc.Features {
		bb := feat.Geometry.BBox()
		bboxPts = append(bboxPts, bb.SW, bb.NE)
	}

	log.Println("determining which tiles need to be generated")
	var tc []tile.ID
	for _, zoomlevel := range zoomlevels {
		tc = append(tc, tile.Coverage(spatial.Line(bboxPts).BBox(), zoomlevel)...)
	}

	var fts spatial.Filterable
	if len(fc.Features)*len(tc) > indexThreshold {
		log.Println("building index...")
		fts = spatial.NewRTreeCollection(fc.Features...)
		log.Println("index complete")
	} else {
		fts = &fc
	}

	log.Printf("starting to generate %d tiles...", len(tc))
	dtw := diskTileWriter{basedir: *target, compressTiles: *compressTiles}
	dlm := defaultLayerMapper{defaultLayer: *defaultLayer}

	var (
		wg sync.WaitGroup
		ws = workerSlices(tc, *workersNumber)
		pb = progressbar.NewBar(len(tc), len(ws))
	)
	for wrk := 0; wrk < len(ws); wrk++ {
		wg.Add(1)
		go func(i int) {
			generateTiles(ws[i], fts, &dtw, &dlm, pb)
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
	basedir       string
	compressTiles bool
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

	if tw.compressTiles {
		_, err = gzip.NewWriter(tf).Write(buf)
	} else {
		_, err = tf.Write(buf)
	}
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

func generateTiles(tIDs []tile.ID, features spatial.Filterable, tw tileWriter, lm layerMapper, pb chan<- struct{}) {
	for _, tID := range tIDs {
		// if !*quiet {
		// 	log.Printf("Generating %s", tID)
		// }
		var (
			layers = map[string][]spatial.Feature{}
			ln     string
		)
		tileClipBBox := tID.BBox()

		for _, feat := range features.Filter(tileClipBBox) {
			sf := tile.Resolution(tID.Z, 4096) * 20
			gm := feat.Geometry.Simplify(sf)
			for _, geom := range gm.ClipToBBox(tileClipBBox) {
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
			pb <- struct{}{}
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
		pb <- struct{}{}
	}
	time.Sleep(100 * time.Millisecond)
}

func anyFeatures(layers map[string][]spatial.Feature) bool {
	for _, ly := range layers {
		if len(ly) > 0 {
			return true
		}
	}
	return false
}
