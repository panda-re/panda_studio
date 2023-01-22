package db

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	config "github.com/panda-re/panda_studio/internal/configuration"
)

var s3Client *minio.Client

func GetS3Client(ctx context.Context) (*minio.Client, error) {
	if s3Client != nil {
		return s3Client, nil
	}

	opts := config.GetConfig().S3

	client, err := minio.New(opts.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(opts.AccessKey, opts.SecretKey, ""),
		Secure: opts.SslEnabled,
	})
	if err != nil {
		return nil, err
	}

	s3Client = client

	return client, err
}