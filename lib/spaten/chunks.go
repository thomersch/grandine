package spaten

import (
	"io"
	"sync"

	"github.com/thomersch/grandine/lib/spatial"
)

type Chunks struct {
	endReached bool
	reader     io.Reader
	// Parallel reading of a file is not allowed, could be theoretically improved by reading from
	// stream and passing the buffer into the decoder, but this needs underlying changes.
	readerMtx sync.Mutex
}

func (c *Chunks) Next() bool {
	return !c.endReached
}

func (c *Chunks) Scan(fc *spatial.FeatureCollection) error {
	c.readerMtx.Lock()
	defer c.readerMtx.Unlock()
	err := readBlock(c.reader, fc)
	if err == io.EOF {
		c.endReached = true
	}
	return nil
}
