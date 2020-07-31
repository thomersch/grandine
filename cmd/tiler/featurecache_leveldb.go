package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/jmhodges/levigo"

	"github.com/thomersch/grandine/lib/spaten"
	"github.com/thomersch/grandine/lib/spatial"
	"github.com/thomersch/grandine/lib/tile"
)

type LevelDBCache struct {
	Zoomlevels []int
	CacheSize  int

	db     *levigo.DB
	dbpath string

	cache *FeatureMap

	count int
	bbox  *spatial.BBox
}

func NewLevelDBCache(zl []int) (*LevelDBCache, error) {
	ldbopt := levigo.NewOptions()
	ldbopt.SetCreateIfMissing(true)
	ldbopt.SetWriteBufferSize(100000000)
	ldbopt.SetCompression(levigo.NoCompression)

	dbpath, err := ioutil.TempDir("", "tiler-leveldb")
	if err != nil {
		return nil, err
	}

	leveldb, err := levigo.Open(dbpath, ldbopt)
	if err != nil {
		return nil, err
	}

	return &LevelDBCache{
		Zoomlevels: zl,
		CacheSize:  100000,
		db:         leveldb,
		dbpath:     dbpath,
		cache:      NewFeatureMap(zl),
	}, nil
}

func (ldb *LevelDBCache) Close() error {
	ldb.flush()
	ldb.db.Close()
	os.RemoveAll(ldb.dbpath)
	return nil
}

func (ldb *LevelDBCache) flush() error {
	var c spaten.Codec
	for tid, fts := range ldb.cache.Dump() {
		tidb := []byte(tid.String())

		rbuf, err := ldb.db.Get(levigo.NewReadOptions(), tidb)
		if err != nil {
			return err
		}

		if len(rbuf) == 0 {
			var wbuf bytes.Buffer
			err := c.Encode(&wbuf, &spatial.FeatureCollection{Features: fts})
			if err != nil {
				return err
			}
			err = ldb.db.Put(levigo.NewWriteOptions(), tidb, wbuf.Bytes())
			if err != nil {
				return err
			}
		} else {
			var (
				buf  = bytes.NewReader(rbuf)
				coll spatial.FeatureCollection
			)
			err := c.Decode(buf, &coll)
			if err != nil {
				return err
			}
			coll.Features = append(coll.Features, fts...)
			var ob bytes.Buffer
			err = c.Encode(&ob, &coll)
			if err != nil {
				return err
			}
			err = ldb.db.Put(levigo.NewWriteOptions(), tidb, ob.Bytes())
			if err != nil {
				return err
			}
		}
	}
	ldb.cache.Clear()
	return nil
}

func (ldb *LevelDBCache) AddFeature(ft spatial.Feature) {
	ldb.count++

	if ldb.count%ldb.CacheSize == 0 {
		err := ldb.flush()
		if err != nil {
			panic(err)
		}
		log.Printf("Written %v features to disk", ldb.count)
	}

	ldb.cache.AddFeature(ft)

	if ldb.bbox == nil {
		var bb = ft.Geometry.BBox()
		ldb.bbox = &bb
	} else {
		ldb.bbox.ExtendWith(ft.Geometry.BBox())
	}
}

func (ldb *LevelDBCache) GetFeatures(tid tile.ID) []spatial.Feature {
	var (
		c  spaten.Codec
		fc spatial.FeatureCollection
	)
	tidb := []byte(tid.String())

	buf, err := ldb.db.Get(levigo.NewReadOptions(), tidb)
	if err != nil {
		panic(err)
	}
	if len(buf) == 0 {
		return nil
	}
	var b = bytes.NewReader(buf)
	err = c.Decode(b, &fc)
	if err != nil {
		panic(err)
	}
	return fc.Features
}

func (ldb *LevelDBCache) BBox() spatial.BBox {
	return *ldb.bbox
}

func (ldb *LevelDBCache) Count() int {
	return ldb.count
}
