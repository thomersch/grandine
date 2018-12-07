package geojsonseq

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomersch/grandine/lib/spatial"
)

func TestChunkedDecode(t *testing.T) {
	f, err := os.Open("testdata/example.geojsonseq")
	assert.Nil(t, err)
	defer f.Close()

	c := Codec{}
	chunks, err := c.ChunkedDecode(f)
	assert.Nil(t, err)

	var fcoll spatial.FeatureCollection
	for chunks.Next() {
		err := chunks.Scan(&fcoll)
		assert.Nil(t, err)
	}
	assert.Len(t, fcoll.Features, 10)
}

type gjProducer struct {
	Data     []byte
	NRecords int
	pos      int
}

func (p *gjProducer) Read(buf []byte) (int, error) {
	var read int

	if p.pos == p.NRecords*len(p.Data) {
		return 0, io.EOF
	}
	for i := 0; i < len(buf); i++ {
		buf[i] = p.Data[p.pos%len(p.Data)]
		p.pos++
		read++
		if p.pos == p.NRecords*len(p.Data) {
			break
		}
	}
	return read, nil
}

func BenchmarkChunkedDecode(b *testing.B) {
	var g = gjProducer{
		NRecords: b.N,
		Data: []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[8.9026586,53.100125]},"properties":{}}
`),
	}

	c := Codec{}

	b.ReportAllocs()
	b.ResetTimer()

	chunks, err := c.ChunkedDecode(&g)
	assert.Nil(b, err)

	var fcoll spatial.FeatureCollection

	for chunks.Next() {
		err := chunks.Scan(&fcoll)
		assert.Nil(b, err)
	}
	assert.Len(b, fcoll.Features, b.N)
}
