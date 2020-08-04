package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"

	humanize "github.com/dustin/go-humanize"

	"github.com/thomersch/grandine/lib/mvt"
	"github.com/thomersch/grandine/lib/progressbar"
	"github.com/thomersch/grandine/lib/spaten"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

type zmLvl []int

func (zm *zmLvl) String() string {
	return fmt.Sprintf("%d", *zm)
}

func (zm *zmLvl) Set(value string) error {
	for _, s := range strings.Split(value, ",") {
		zs := strings.TrimSpace(s)
		if zs == "" {
			continue
		}
		v, err := strconv.Atoi(zs)
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
	var (
		sourceStdIn bool
		tileCodec   tile.Codec
	)
	source := flag.String("in", "", "file to read from, supported format: spaten")
	target := flag.String("out", "tiles", "path where the tiles will be written")
	defaultLayer := flag.Bool("default-layer", true, "if no layer name is specified in the feature, whether it will be put into a default layer")
	workersNumber := flag.Int("workers", runtime.GOMAXPROCS(0), "number of workers")
	cpuProfile := flag.String("cpuprof", "", "writes CPU profiling data into a file")
	geojsonCodec := flag.Bool("geojson", false, "encode tiles into geojson instead of MVT, for debugging purposes")
	compressTiles := flag.Bool("compress", false, "compress tiles with gzip")
	cacheStrategy := flag.String("cache", "leveldb", fmt.Sprintf("cache strategy, possible values: %v", availableCaches()))
	quiet = flag.Bool("q", false, "argument to use if program should be run in quiet mode with reduced logging")

	flag.Var(&zoomlevels, "zoom", "one or more zoom levels (comma separated) of which the tiles will be rendered")
	flag.Parse()

	if len(*source) == 0 {
		log.Print("grandine-tiler: Reading from stdin.")
		sourceStdIn = true
	}

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

	if *geojsonCodec {
		tileCodec = &tile.GeoJSONCodec{}
	} else {
		tileCodec = &mvt.Codec{}
	}

	var (
		f   io.Reader
		err error
	)

	if !sourceStdIn {
		f, err = os.Open(*source)
		if err != nil {
			log.Fatal(err)
		}
		defer f.(io.Closer).Close()
	} else {
		f = os.Stdin
	}

	var tw tileWriter

	if strings.HasPrefix(*target, S3ProtoPrefix) {
		u, err := url.Parse(*target)
		if err != nil {
			log.Fatalf("Assumed target path is S3 URL, but encoutered error while parsing: %s", err)
		}

		s3key := os.Getenv("S3KEY")
		s3secret := os.Getenv("S3SECRET")
		tw, err = NewS3TileWriter(u.Host, strings.TrimPrefix(u.Path, "/"), s3key, s3secret)
		if err != nil {
			log.Fatalf("Could not create S3 Client: %s", err)
		}
	} else {
		err = os.MkdirAll(*target, 0777)
		if err != nil {
			log.Fatal(err)
		}
		tw = &diskTileWriter{basedir: *target, compressTiles: *compressTiles}
	}

	log.Println("Preparing feature table...")

	cinit, ok := caches[*cacheStrategy]
	if !ok {
		log.Fatalf("invalid cache strategy name '%s', available: %v", *cacheStrategy, availableCaches())
	} else if !*quiet {
		log.Printf("Using %s cache", *cacheStrategy)
	}
	ft, err := cinit(zoomlevels)
	if err != nil {
		log.Fatal(err)
	}

	defer func(ft FeatureCache) {
		clsr, ok := ft.(io.Closer)
		if ok {
			log.Println("Cleaning up cache.")
			clsr.Close()
		}
	}(ft)
	showMemStats()

	log.Println("Parsing input...")

	var codec spaten.Codec
	cd, err := codec.ChunkedDecode(f)
	if err != nil {
		log.Fatalf("Could not read inscoming file: %v", err)
	}

	var fc spatial.FeatureCollection
	for cd.Next() {
		cd.Scan(&fc)
		for _, feat := range fc.Features {
			ft.AddFeature(feat)
		}
		fc.Reset()
	}
	log.Printf("%v feature are in-cache", ft.Count())
	showMemStats()

	log.Println("Determining which tiles need to be generated")
	var tc []tile.ID
	for _, zoomlevel := range zoomlevels {
		tc = append(tc, tile.Coverage(ft.BBox(), zoomlevel)...)
	}

	log.Printf("Starting to generate %d tiles...", len(tc))

	dlm := defaultLayerMapper{defaultLayer: *defaultLayer}

	shuffleWork(tc) // randomize order for better worker saturation
	var (
		wg       sync.WaitGroup
		ws       = workerSlices(tc, *workersNumber)
		pb, done = progressbar.NewBar(len(tc), len(ws)) // TODO: respect quiet setting
	)
	for wrk := 0; wrk < len(ws); wrk++ {
		wg.Add(1)
		go func(i int) {
			generateTiles(ws[i], ft, tw, tileCodec, &dlm, pb)
			wg.Done()
		}(wrk)
	}
	wg.Wait()
	done()

	showMemStats()
	log.Println("Done.")
}

func renderable(props map[string]interface{}, zl int) bool {
	var propInt = func(props map[string]interface{}, name string, defaultVal int) int {
		v, ok := props[name]
		if !ok {
			return defaultVal
		}
		i, ok := v.(int)
		if ok {
			return i
		}
		f, ok := v.(float64)
		if ok {
			return int(f)
		}
		log.Printf("%v is neither int nor float: %v", props, v)
		return defaultVal
	}
	return (zl >= propInt(props, "@zoom:min", 0)) && (zl <= propInt(props, "@zoom:max", 99))
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

func (tw *diskTileWriter) WriteTile(tID tile.ID, buf []byte, ext string) error {
	err := os.MkdirAll(filepath.Join(tw.basedir, strconv.Itoa(tID.Z), strconv.Itoa(tID.X)), 0777)
	if err != nil {
		return err
	}
	tf, err := os.Create(filepath.Join(tw.basedir, strconv.Itoa(tID.Z), strconv.Itoa(tID.X), strconv.Itoa(tID.Y)+"."+ext))
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
	WriteTile(tile.ID, []byte, string) error
}

func generateTiles(tIDs []tile.ID, fts FeatureCache, tw tileWriter, encoder tile.Codec, lm layerMapper, pb chan<- struct{}) {
	for _, tID := range tIDs {
		var (
			layers = map[string][]spatial.Feature{}
			ln     string
		)

		for _, feat := range fts.GetFeatures(tID) {
			ln = lm.LayerName(feat.Props)
			if len(ln) != 0 {
				if _, ok := layers[ln]; !ok {
					layers[ln] = []spatial.Feature{feat}
				} else {
					layers[ln] = append(layers[ln], feat)
				}
			}
		}

		pb <- struct{}{}

		if !anyFeatures(layers) {
			continue
		}
		buf, err := encoder.EncodeTile(layers, tID)
		if err != nil {
			log.Fatal(err)
		}

		err = tw.WriteTile(tID, buf, encoder.Extension())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func anyFeatures(layers map[string][]spatial.Feature) bool {
	for i := range layers {
		if len(layers[i]) > 0 {
			return true
		}
	}
	return false
}

func pow(x, y int) int {
	var res = 1
	for i := 1; i <= y; i++ {
		res *= x
	}
	return res
}

func showMemStats() {
	if *quiet {
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Memory in use: %s", humanize.Bytes(m.Alloc))
}
