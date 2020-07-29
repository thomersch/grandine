// +build !go1.13

package main

import "github.com/thomersch/grandine/lib/tile"

const S3ProtoPrefix = "s3://"

type S3TileWriter struct{}

func (s3 *S3TileWriter) WriteTile(tID tile.ID, buf []byte, ext string) error {
	panic("S3TileWriter needs at least Go 1.13")
}

func NewS3TileWriter(endpoint, bucket, accessKeyID, secretAccessKey string) (*S3TileWriter, error) {
	panic("S3TileWriter needs at least Go 1.13")
}
