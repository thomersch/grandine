package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/thomersch/grandine/lib/csv"
	"github.com/thomersch/grandine/lib/geojson"
	"github.com/thomersch/grandine/lib/mapping"
	"github.com/thomersch/grandine/lib/spaten"
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

func main() {
	var (
		infiles filelist
		conds   []mapping.Condition
	)
	dest := flag.String("out", "geo.spaten", "")
	mapFilePath := flag.String("mapping", "", "Path to mapping file which will be used to transform data.")
	csvLatColumn := flag.Int("csv-lat", 1, "If parsing CSV, which column contains the Latitude. Zero-indexed.")
	csvLonColumn := flag.Int("csv-lon", 2, "If parsing CSV, which column contains the Longitude. Zero-indexed.")
	csvDelimiter := flag.String("csv-delim", ",", "If parsing CSV, what is the delimiter between values")
	flag.Var(&infiles, "in", "infile(s)")
	flag.Parse()

	if len(*csvDelimiter) > 1 {
		log.Fatal("CSV Delimiter: only single character delimiters are allowed")
	}

	if len(*mapFilePath) != 0 {
		mf, err := os.Open(*mapFilePath)
		if err != nil {
			log.Fatal(err)
		}
		conds, err = mapping.ParseMapping(mf)
		if err != nil {
			log.Fatal(err)
		}

		if len(conds) > 0 {
			log.Printf("input file(s) will be filtered using %v conditions", len(conds))
		}
	}

	availableCodecs := []spatial.Codec{
		&geojson.Codec{},
		&spaten.Codec{},
		&csv.Codec{
			LatCol: *csvLatColumn,
			LonCol: *csvLonColumn,
			Delim:  rune((*csvDelimiter)[0]),
		},
	}

	enc, err := guessCodec(*dest, availableCodecs)
	if err != nil {
		log.Fatalf("file type of %s is not supported (please check for correct file extension)", *dest)
	}
	encoder, ok := enc.(spatial.Encoder)
	if !ok {
		log.Fatalf("%v codec does not support writing", enc)
	}

	out, err := os.Create(*dest)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	var fc spatial.FeatureCollection
	for _, infileName := range infiles {
		dec, err := guessCodec(infileName, availableCodecs)
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

		var finished func() error
		switch d := dec.(type) {
		case spatial.ChunkedDecoder:
			chunks, err := d.ChunkedDecode(r)
			if err != nil {
				log.Fatalf("could not decode %v: %v", infileName, err)
			}
			for chunks.Next() {
				chunks.Scan(&fc)
				finished, err = write(out, &fc, encoder, conds)
				if err != nil {
					log.Fatal(err)
				}
				fc.Features = []spatial.Feature{}
			}
			err = finished()
			if err != nil {
				log.Fatal(err)
			}
		case spatial.Decoder:
			err = decoder.Decode(r, &fc)
			if err != nil {
				log.Fatalf("could not decode %v: %v", infileName, err)
			}
			finished, err = write(out, &fc, encoder, conds)
			if err != nil {
				log.Fatal(err)
			}
			finished()
		}
	}
}

var featBuf []spatial.FeatureCollection // TODO: this is not optimal, needs better wrapping

func write(w io.Writer, fs *spatial.FeatureCollection, enc spatial.Encoder, conds []mapping.Condition) (flush func() error, err error) {
	if len(conds) > 0 {
		var filtered []spatial.Feature
		for _, ft := range fs.Features {
			for _, cond := range conds {
				if cond.Matches(ft.Props) {
					nft := ft
					nft.Props = cond.Map(ft.Props)
					filtered = append(filtered, nft)
				}
			}
		}
		fs.Features = filtered
	}

	if e, ok := enc.(spatial.ChunkedEncoder); ok {
		err = e.EncodeChunk(w, fs)
		if err != nil {
			return func() error { return nil }, err
		}
		return func() error { return e.Close() }, nil
	}

	featBuf = append(featBuf, *fs)
	return func() error {
		var flat spatial.FeatureCollection
		for _, ftc := range featBuf {
			flat.Features = append(flat.Features, ftc.Features...)
			flat.SRID = ftc.SRID
		}
		return enc.Encode(w, &flat)
	}, nil
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
