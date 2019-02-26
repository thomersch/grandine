package geojsonseq

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/thomersch/grandine/lib/geojson"
	"github.com/thomersch/grandine/lib/spatial"
)

const resourceSep = byte('\x1E')

type Codec struct{}

func (c *Codec) Decode(io.Reader, *spatial.FeatureCollection) error {
	panic("not implemented yet, please use ChunkedDecode")
}

func (c *Codec) ChunkedDecode(r io.Reader) (spatial.Chunks, error) {
	var rs = make([]byte, 1)
	_, err := r.Read(rs)
	if rs[0] != resourceSep || err != nil {
		return nil, fmt.Errorf("a geojson sequence must start with a resource separator, got: %v", rs[0])
	}
	return &chunk{reader: bufio.NewReader(r)}, nil
}

func (c *Codec) Extensions() []string {
	return []string{"geojsonseq"}
}

type chunk struct {
	endReached bool
	reader     *bufio.Reader
	buf        []byte
	err        error
}

func (ch *chunk) Next() bool {
	if ch.endReached {
		return false
	}
	var err error
	ch.buf, err = ch.reader.ReadBytes(resourceSep)
	if err == io.EOF {
		ch.endReached = true
	} else if err != nil {
		ch.err = err // we cannot return the value here, so it will pop up on the next scan
		return true
	}
	ch.buf = ch.buf[:len(ch.buf)-1]
	return true
}

func (ch *chunk) Scan(fc *spatial.FeatureCollection) error {
	if ch.err != nil {
		return ch.err
	}
	var fts geojson.FeatList
	err := json.Unmarshal(append(append([]byte(`[`), ch.buf...), ']'), &fts)
	if err != nil {
		return err
	}
	fc.Features = append(fc.Features, fts...)
	return nil
}
