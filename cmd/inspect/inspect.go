package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/thomersch/grandine/lib/spaten"
	"github.com/thomersch/grandine/lib/spatial"
)

func colorKV(s string) string {
	if strings.HasPrefix(s, "@") {
		return "\033[37m" + s + "\033[0m"
	}
	return s
}

func prettyPrint(m map[string]interface{}) string {
	var (
		i  int
		ol = make([]string, len(m))
		o  string
	)

	for k := range m {
		ol[i] = k
		i++
	}
	sort.Strings(ol)
	for _, k := range ol {
		o += colorKV(fmt.Sprintf("%v=%v ", k, m[k]))
	}
	return o
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s filepath\n", os.Args[0])
		os.Exit(1)
	}
	filepath := os.Args[len(os.Args)-1]
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Could not open %v", filepath)
	}
	defer f.Close()

	pp := exec.Command("less", "-R")
	pp.Stdout = os.Stdout
	stdin, err := pp.StdinPipe()
	if err != nil {
		log.Fatalf("Could not open pipe for pager: %v", err)
	}
	pp.Start()

	go func() {
		defer stdin.Close()

		var (
			codec spaten.Codec
			fc    = spatial.NewFeatureCollection()
		)
		hd, err := spaten.ReadFileHeader(f)
		if err != nil {
			io.WriteString(stdin, fmt.Sprintf("Invalid Header: %v", err))
			return
		}
		fmt.Fprintf(stdin, "Spaten file, Version %v\n", hd.Version)
		_, err = f.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}

		chunks, err := codec.ChunkedDecode(f)
		if err != nil {
			log.Fatalf("could not decode: %v", err)
		}
		for chunks.Next() {
			err = chunks.Scan(fc)
			if err != nil {
				log.Fatal(err)
			}
			for _, ft := range fc.Features {
				fmt.Fprintf(stdin, "\033[34m%v\033[0m ", ft.Geometry.Typ())
				fmt.Fprintf(stdin, "%v", prettyPrint(ft.Props))
				fmt.Fprintf(stdin, "\033[31m%v\033[0m\n", ft.Geometry)
			}
			fc.Features = []spatial.Feature{}
		}
	}()
	pp.Wait()
}
