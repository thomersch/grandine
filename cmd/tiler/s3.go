package main

import (
	"bytes"
	"context"

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
		// Secure: true,
	})
	s3tw.client = s3Client{client, bucket}
	if err != nil {
		return s3tw, err
	}

	return s3tw, nil
}

func (s3 *S3TileWriter) WriteTile(tID tile.ID, buf []byte, ext string) error {
	r := bytes.NewReader(buf)

	_, err := s3.client.PutObject(context.Background(), s3.client.bucket, tID.String()+"."+ext, r, int64(r.Len()), minio.PutObjectOptions{})
	return err
}
