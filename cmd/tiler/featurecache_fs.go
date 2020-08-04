package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/thomersch/grandine/lib/spaten"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

type FileSystemCache struct {
	Zoomlevels []int
	CacheSize  int

	basePath string
	fp       map[tile.ID]*os.File

	cache *FeatureMap

	count          int
	lastCheckpoint time.Time

	bbox *spatial.BBox
}

func NewFileSystemCache(zl []int) (*FileSystemCache, error) {
	basePath, err := ioutil.TempDir("", "tiler-fscache")
	if err != nil {
		return nil, err
	}

	return &FileSystemCache{
		Zoomlevels: zl,
		CacheSize:  1000000,
		basePath:   basePath,
		fp:         make(map[tile.ID]*os.File),

		cache:          NewFeatureMap(zl),
		lastCheckpoint: time.Now(),
	}, nil
}

func (fsc *FileSystemCache) Close() error {
	os.RemoveAll(fsc.basePath)
	return nil
}

func (fsc *FileSystemCache) fpath(tid tile.ID) string {
	return path.Join(fsc.basePath, strconv.Itoa(tid.Z)+"-"+strconv.Itoa(tid.X)+"-"+strconv.Itoa(tid.Y))
}

func (fsc *FileSystemCache) flush() error {
	var (
		c   spaten.Codec
		err error
	)
	for tid, fts := range fsc.cache.Dump() {
		fd, ok := fsc.fp[tid]
		if !ok {
			var wbuf bytes.Buffer
			err := c.Encode(&wbuf, &spatial.FeatureCollection{Features: fts})
			if err != nil {
				return err
			}

			f, err := os.Create(fsc.fpath(tid))
			if err != nil {
				return err
			}
			fsc.fp[tid] = f
			wbuf.WriteTo(f)
		} else {
			err = spaten.WriteBlock(fd, fts, nil)
			if err != nil {
				return err
			}
		}
	}
	fsc.cache.Clear()
	return nil
}

func (fsc *FileSystemCache) AddFeature(ft spatial.Feature) {
	fsc.count++

	if fsc.count%fsc.CacheSize == 0 {
		showMemStats()
		err := fsc.flush()
		if err != nil {
			panic(err)
		}

		log.Printf("Written %v features to disk (%.0f/s)", fsc.count, float64(fsc.CacheSize)/time.Since(fsc.lastCheckpoint).Seconds())
		fsc.lastCheckpoint = time.Now()
	}

	fsc.cache.AddFeature(ft)

	if fsc.bbox == nil {
		var bb = ft.Geometry.BBox()
		fsc.bbox = &bb
	} else {
		fsc.bbox.ExtendWith(ft.Geometry.BBox())
	}
}

func (fsc *FileSystemCache) GetFeatures(tid tile.ID) []spatial.Feature {
	var (
		c  spaten.Codec
		fc spatial.FeatureCollection
	)

	fp, ok := fsc.fp[tid]
	if !ok {
		return nil
	}

	_, err := fp.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	err = c.Decode(fp, &fc)
	if err != nil {
		panic(err)
	}
	return fc.Features
}

func (fsc *FileSystemCache) BBox() spatial.BBox {
	return *fsc.bbox
}

func (fsc *FileSystemCache) Count() int {
	return fsc.count
}
