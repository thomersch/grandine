package cugdf

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/thomersch/grandine/converter/fileformat"
	"github.com/thomersch/grandine/lib/spatial"

	"github.com/golang/protobuf/proto"
)

const (
	cookie  = "ABCD"
	version = 0
)

var encoding = binary.LittleEndian

func WriteFileHeader(w io.Writer) error {
	const headerSize = 8

	buf := make([]byte, headerSize)
	buf = append([]byte(cookie), buf[:4]...)
	binary.LittleEndian.PutUint32(buf[4:], version)

	n, err := w.Write(buf)
	if n != headerSize {
		return io.EOF
	}
	return err
}

func WriteBlock(w io.Writer, fs []spatial.Feature) error {
	blockBody := &fileformat.Body{}
	for _, f := range fs {
		var (
			tags []*fileformat.Tag
		)
		for k, v := range f.Properties() {
			val, typ, err := fileformat.ValueType(v)
			if err != nil {
				return err
			}
			tags = append(tags, &fileformat.Tag{
				Key:   k,
				Value: val,
				Type:  typ,
			})
		}

		wkbBuf, err := f.MarshalWKB()
		if err != nil {
			return err
		}

		blockBody.Feature = append(blockBody.Feature, &fileformat.Feature{
			Geom: wkbBuf,
			Tags: tags,
		})
	}
	bodyBuf, err := proto.Marshal(blockBody)
	if err != nil {
		log.Fatal(err)
	}

	// Body Length (fill later)
	binary.Write(w, binary.LittleEndian, uint32(len(bodyBuf)))
	// Flags
	binary.Write(w, binary.LittleEndian, uint16(0))
	// Compression
	binary.Write(w, binary.LittleEndian, uint8(0))
	// Message Type
	binary.Write(w, binary.LittleEndian, uint8(0))

	w.Write(bodyBuf)
	return nil
}
