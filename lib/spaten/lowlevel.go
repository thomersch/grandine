package spaten

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/thomersch/grandine/lib/spaten/fileformat"
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
		return hd, fmt.Errorf("could not read file header cookie: %s", err)
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

// WriteBlock writes a block of spatial data (note that every valid Spaten file needs a file header in front).
// meta may be nil, if you don't wish to add any block meta.
func WriteBlock(w io.Writer, fs []spatial.Feature, meta map[string]interface{}) error {
	blockBody := &fileformat.Body{}
	props, err := propertiesToTags(meta)
	if err != nil {
		return err
	}
	blockBody.Meta = &fileformat.Meta{
		Tags: props,
	}

	for _, f := range fs {
		var nf fileformat.Feature
		nf.Tags, err = propertiesToTags(f.Properties())
		if err != nil {
			return err
		}

		// TODO: make encoder configurable
		nf.Geom, err = f.MarshalWKB()
		if err != nil {
			return err
		}

		blockBody.Feature = append(blockBody.Feature, &nf)
	}
	bodyBuf, err := proto.Marshal(blockBody)
	if err != nil {
		log.Fatal(err)
	}

	blockHeaderBuf := make([]byte, 8)
	// Body Length
	binary.LittleEndian.PutUint32(blockHeaderBuf[:4], uint32(len(bodyBuf)))
	// Flags
	binary.LittleEndian.PutUint16(blockHeaderBuf[4:6], 0)
	// Compression
	blockHeaderBuf[6] = 0
	// Message Type
	blockHeaderBuf[7] = 0

	w.Write(append(blockHeaderBuf, bodyBuf...))
	return nil
}

func propertiesToTags(props map[string]interface{}) ([]*fileformat.Tag, error) {
	var tags []*fileformat.Tag
	if props == nil {
		return tags, nil
	}
	for k, v := range props {
		val, typ, err := fileformat.ValueType(v)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &fileformat.Tag{
			Key:   k,
			Value: val,
			Type:  typ,
		})
	}
	return tags, nil
}

func readBlock(r io.Reader, fs *spatial.FeatureCollection) error {
	var (
		blockLength uint32
		flags       uint16
		compression uint8
		messageType uint8
	)
	if err := binary.Read(r, binary.LittleEndian, &blockLength); err != nil {
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
	if n, err := r.Read(buf); err != nil {
		return err
	} else if n != int(blockLength) {
		return fmt.Errorf("incomplete block: expected %v bytes, %v available", blockLength, n)
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
	return nil
}

// ReadBlocks is a function for reading all features from a file at once.
func ReadBlocks(r io.Reader, fs *spatial.FeatureCollection) error {
	var err error
	for {
		err = readBlock(r, fs)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
