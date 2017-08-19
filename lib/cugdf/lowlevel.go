package cugdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"

	"github.com/thomersch/grandine/lib/cugdf/fileformat"
	"github.com/thomersch/grandine/lib/spatial"

	"github.com/golang/protobuf/proto"
)

const (
	cookie  = "SPAT"
	version = 0
)

type Header struct {
	Version int
}

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

func ReadFileHeader(r io.Reader) (Header, error) {
	var (
		ck   = make([]byte, 4)
		vers uint32
		hd   Header
	)
	if _, err := r.Read(ck); err != nil {
		return hd, err
	}
	if string(ck) != cookie {
		return hd, errors.New("invalid cookie")
	}

	if err := binary.Read(r, binary.LittleEndian, &vers); err != nil {
		return hd, err
	}
	hd.Version = int(vers)
	if vers > version {
		return hd, errors.New("invalid file version")
	}
	return hd, nil
}

func WriteBlock(w io.Writer, fs []spatial.Feature) error {
	blockBody := &fileformat.Body{}
	for _, f := range fs {
		var tags []*fileformat.Tag
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

// ReadBlocks is a high-level interface for reading all features from a file at once.
func ReadBlocks(r io.Reader, fs *spatial.FeatureCollection) error {
	for {
		var (
			blockLength uint32
			flags       uint16
			compression uint8
			messageType uint8
		)
		if err := binary.Read(r, binary.LittleEndian, &blockLength); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &flags); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &compression); err != nil {
			if compression != 0 {
				return errors.New("compression is not supported")
			}
		}
		if err := binary.Read(r, binary.LittleEndian, &messageType); err != nil {
			if messageType != 0 {
				return errors.New("message type is not supported")
			}
		}

		var (
			buf       = make([]byte, blockLength)
			blockBody fileformat.Body
		)
		if _, err := r.Read(buf); err != nil {
			return err
		}
		if err := proto.Unmarshal(buf, &blockBody); err != nil {
			return err
		}
		for _, f := range blockBody.GetFeature() {
			var (
				feature = spatial.Feature{
					Props: map[string]interface{}{},
				}
				geomBuf = bytes.NewBuffer(f.GetGeom())
			)
			err := feature.Geometry.UnmarshalWKB(geomBuf)
			if err != nil {
				return err
			}

			for _, tag := range f.Tags {
				k, v, err := fileformat.KeyValue(tag)
				if err != nil {
					// TODO
					return err
				}
				feature.Props[k] = v
			}
			fs.Features = append(fs.Features, feature)
		}
	}
	return nil
}
