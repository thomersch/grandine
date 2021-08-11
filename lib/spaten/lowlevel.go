package spaten

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

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
		nf, err := PackFeature(f)
		if err != nil {
			return err
		}

		blockBody.Feature = append(blockBody.Feature, &nf)
	}
	bodyBuf, err := proto.Marshal(blockBody)
	if err != nil {
		return err
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

// PackFeature encapusaltes a spatial feature into an encodable Spaten feature.
// This is a low level interface and not guaranteed to be stable.
func PackFeature(f spatial.Feature) (fileformat.Feature, error) {
	var (
		nf  fileformat.Feature
		err error
	)
	nf.Tags, err = propertiesToTags(f.Properties())
	if err != nil {
		return nf, err
	}

	// TODO: make encoder configurable
	nf.Geom, err = f.MarshalWKB()
	if err != nil {
		return nf, err
	}
	return nf, nil
}

var featureBufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 16))
	},
}

// UnpackFeature unpacks a Spaten feature into a usable spatial feature.
// This is a low level interface and not guaranteed to be stable.
func UnpackFeature(pf *fileformat.Feature) (spatial.Feature, error) {
	var geomBuf = featureBufPool.Get().(*bytes.Buffer)
	geomBuf.Reset()
	geomBuf.Write(pf.GetGeom())
	geom, err := spatial.GeomFromWKB(geomBuf)
	if err != nil {
		return spatial.Feature{}, err
	}
	featureBufPool.Put(geomBuf)
	feature := spatial.Feature{
		Props:    map[string]interface{}{},
		Geometry: geom,
	}

	for _, tag := range pf.Tags {
		k, v, err := fileformat.KeyValue(tag)
		if err != nil {
			// TODO
			return feature, err
		}
		feature.Props[k] = v
	}
	return feature, nil
}

func propertiesToTags(props map[string]interface{}) ([]*fileformat.Tag, error) {
	var tags []*fileformat.Tag = make([]*fileformat.Tag, 0, len(props))
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

type blockHeader struct {
	bodyLen     uint32
	flags       uint16
	compression uint8
	messageType uint8
}

var blockBodyPool = sync.Pool{
	New: func() interface{} {
		return &fileformat.Body{}
	},
}

func readBlock(r io.Reader, fs *spatial.FeatureCollection) error {
	var hd blockHeader

	headerBuf := make([]byte, 8)
	n, err := r.Read(headerBuf)
	if n == 0 {
		return io.EOF
	}
	if err != nil {
		return fmt.Errorf("could not read block header: %v", err)
	}

	hd.bodyLen = binary.LittleEndian.Uint32(headerBuf[0:4])
	hd.flags = binary.LittleEndian.Uint16(headerBuf[4:6])
	hd.compression = uint8(headerBuf[6])
	if hd.compression != 0 {
		return errors.New("compression is not supported")
	}

	hd.messageType = uint8(headerBuf[7])
	if hd.messageType != 0 {
		return errors.New("message type is not supported")
	}

	var (
		buf       = make([]byte, hd.bodyLen)
		blockBody = blockBodyPool.Get().(*fileformat.Body)
	)
	blockBody.Reset()
	n, err = io.ReadFull(r, buf)
	if n != int(hd.bodyLen) {
		return fmt.Errorf("incomplete block: expected %v bytes, %v available", hd.bodyLen, n)
	}
	if err != nil {
		return err
	}
	if err := blockBody.Unmarshal(buf); err != nil {
		return err
	}
	if len(fs.Features) == 0 {
		// only prealloc if empty, so no user data gets truncated
		fs.Features = make([]spatial.Feature, 0, len(blockBody.GetFeature()))
	}
	for _, f := range blockBody.GetFeature() {
		feature, err := UnpackFeature(f)
		if err != nil {
			return err
		}
		fs.Features = append(fs.Features, feature)
	}
	blockBodyPool.Put(blockBody)
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
