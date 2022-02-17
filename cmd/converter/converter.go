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
	"github.com/thomersch/grandine/lib/geojsonseq"
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
	dest := flag.String("out", "", "")
	mapFilePath := flag.String("mapping", "", "Path to mapping file which will be used to transform data.")
	csvLatColumn := flag.Int("csv-lat", 1, "If parsing CSV, which column contains the Latitude. Zero-indexed.")
	csvLonColumn := flag.Int("csv-lon", 2, "If parsing CSV, which column contains the Longitude. Zero-indexed.")
	csvDelimiter := flag.String("csv-delim", ",", "If parsing CSV, what is the delimiter between values")
	inCodecName := flag.String("in-codec", "spaten", "Specify codec for in-files. Only used for read from stdin.")
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
		&geojsonseq.Codec{},
	}

	// Determining which codec we will be using for the output.
	var (
		enc interface{}
		err error
	)
	if len(*dest) == 0 {
		enc = &spaten.Codec{}
	} else {
		enc, err = guessCodec(*dest, availableCodecs)
		if err != nil {
			log.Fatalf("file type of %s is not supported (please check for correct file extension)", *dest)
		}
	}
	encoder, ok := enc.(spatial.Encoder)
	if !ok {
		log.Fatalf("%T codec does not support writing", enc)
	}

	// Determine whether we're writing to a stream or file.
	var out io.WriteCloser
	if len(*dest) == 0 {
		out = os.Stdout
	} else {
		out, err = os.Create(*dest)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer out.Close()

	var (
		fc       spatial.FeatureCollection
		finished func() error
	)

	if len(infiles) == 0 {
		log.Println("No input files specified. Reading from stdin.")
		incodec, err := guessCodec("."+*inCodecName, availableCodecs)
		if err != nil {
			log.Fatalf("could not use incodec: %v", err)
		}
		inc, ok := incodec.(spatial.ChunkedDecoder)
		if !ok {
			log.Fatal("codec cannot be used for decoding")
		}
		icd, err := inc.ChunkedDecode(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		for icd.Next() {
			err = icd.Scan(&fc)
			if err != nil {
				log.Fatal(err)
			}
			finished, err = write(out, &fc, encoder, conds)
			fc.Reset()
		}
		finished()
	}

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
					filtered = append(filtered, cond.Transform(ft)...)
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
		return func() error { return e.Close(w) }, nil
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
