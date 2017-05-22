package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/thomersch/grandine/lib/cugdf"
	"github.com/thomersch/grandine/lib/spatial"
)

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
		_, err = cugdf.ReadFileHeader(f)
		if err != nil {
			log.Fatal(err)
		}
		var (
			fcoll spatial.FeatureCollection
		)
		blks, err := cugdf.ReadBlocks(f)
		if err != nil {
			log.Fatal(err)
		}

		for _, ft := range blks {
			fcoll.Features = append(fcoll.Features, ft)
		}
		buf, err := fcoll.MarshalJSON()
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
