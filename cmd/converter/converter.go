package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thomersch/grandine/lib/csv"
	"github.com/thomersch/grandine/lib/cugdf"
	"github.com/thomersch/grandine/lib/cugdf/fileformat"
	"github.com/thomersch/grandine/lib/geojson"
	"github.com/thomersch/grandine/lib/spatial"
)

type filelist []string

func (fl *filelist) String() string {
	return fmt.Sprintf("%s", *fl)
}

func (fl *filelist) Set(value string) error {
	for _, s := range strings.Split(value, ",") {
		*fl = append(*fl, strings.TrimSpace(s))
	}
	return nil
}

var codecs = []spatial.Codec{
	&geojson.Codec{},
	&cugdf.Codec{},
	&csv.Codec{
		//TODO: make configurable via flags
		LatCol: 4,
		LonCol: 5,
		ColPropMap: map[int]csv.TagMapping{
			1:  {"name", fileformat.Tag_STRING},
			14: {"population", fileformat.Tag_INT},
		},
	},
}

func main() {
	var infiles filelist
	dest := flag.String("out", "geo.spaten", "")
	flag.Var(&infiles, "in", "infile(s)")
	flag.Parse()

	enc, err := guessCodec(*dest, codecs)
	if err != nil {
		log.Fatalf("file type of %s is not supported (please check for correct file extension)", *dest)
	}
	encoder, ok := enc.(spatial.Encoder)
	if !ok {
		log.Fatalf("%v codec does not support writing", enc)
	}

	var fc spatial.FeatureCollection
	for _, infileName := range infiles {
		dec, err := guessCodec(infileName, codecs)
		if err != nil {
			log.Fatalf("file type of %s is not supported (please check for correct file extension)", infileName)
		}
		decoder, ok := dec.(spatial.Decoder)
		if !ok {
			log.Fatalf("%T codec does not support reading", decoder)
		}

		r, err := os.Open(infileName)
		if err != nil {
			log.Fatalf("could not open %v for reading: %v", infileName, err)
		}
		defer r.Close()
		err = decoder.Decode(r, &fc)
		if err != nil {
			log.Fatalf("could not decode %v: %v", infileName, err)
		}
	}

	out, err := os.Create(*dest)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// spew.Dump(fc)
	err = encoder.Encode(out, &fc)
	if err != nil {
		log.Fatalf("could not encode %v: %v", *dest, err)
	}
}

func guessCodec(filename string, codecs []spatial.Codec) (spatial.Codec, error) {
	fn := strings.ToLower(filename)
	for _, cd := range codecs {
		for _, ext := range cd.Extensions() {
			if strings.HasSuffix(fn, "."+ext) {
				return cd, nil
			}
		}
	}
	return nil, errors.New("file type not supported")
}
