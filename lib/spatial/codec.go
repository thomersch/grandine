package spatial

import "io"

type Chunks interface {
	Next() bool
	Scan(fc *FeatureCollection) error
}

type Decoder interface {
	Decode(io.Reader, *FeatureCollection) error
}

type ChunkedDecoder interface {
	ChunkedDecode(io.Reader) (Chunks, error)
}

type Encoder interface {
	Encode(io.Writer, *FeatureCollection) error
}

// A Codec needs to be able to tell which file extensions (e.g. "geojson")
// are commonly used to persist files. Moreover a Codec SHOULD either implement
// a Decoder or Encoder.
type Codec interface {
	Extensions() []string
}
