package cugdf

import (
	"io"

	"github.com/thomersch/grandine/lib/spatial"
)

// Marshal is a high-level interface for writing features into a file.
func Marshal(feats []spatial.Feature, w io.Writer) error {
	err := WriteFileHeader(w)
	if err != nil {
		return err
	}

	for _, ftBlk := range geomBlocks(100, feats) {
		err = WriteBlock(w, ftBlk)
		if err != nil {
			return err
		}
	}
	return nil
}

func Unmarshal(r io.Reader) ([]spatial.Feature, error) {
	_, err := ReadFileHeader(r)
	if err != nil {
		return nil, err
	}
	return ReadBlocks(r)
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
