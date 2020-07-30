// +build go1.13

package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thomersch/grandine/lib/tile"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const S3ProtoPrefix = "s3://"

type s3Client struct {
	*minio.Client
	bucket string
}

type S3TileWriter struct {
	client s3Client
}

func NewS3TileWriter(endpoint, bucket, accessKeyID, secretAccessKey string) (*S3TileWriter, error) {
	s3tw := &S3TileWriter{}
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	ok, err := client.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("S3 bucket %s does not exist", bucket)
	}

	s3tw.client = s3Client{client, bucket}

	return s3tw, nil
}

func (s3 *S3TileWriter) WriteTile(tID tile.ID, buf []byte, ext string) error {
	r := bytes.NewReader(buf)

	_, err := s3.client.PutObject(context.Background(), s3.client.bucket, tID.String()+"."+ext, r, int64(r.Len()), minio.PutObjectOptions{})

	return errors.Wrap(err, "S3 upload failed")
}
