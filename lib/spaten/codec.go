package spaten

import (
	"io"

	"github.com/thomersch/grandine/lib/spatial"
)

type Codec struct{}

func (c *Codec) Encode(w io.Writer, fc *spatial.FeatureCollection) error {
	err := WriteFileHeader(w)
	if err != nil {
		return err
	}

	for _, ftBlk := range geomBlocks(100, fc.Features) {
		err = WriteBlock(w, ftBlk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Codec) Decode(r io.Reader, fc *spatial.FeatureCollection) error {
	_, err := ReadFileHeader(r)
	if err != nil {
		return err
	}
	err = ReadBlocks(r, fc)
	return err
}

func (c *Codec) Extensions() []string {
	return []string{"spaten"}
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
