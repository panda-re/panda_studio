package main

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	config "github.com/panda-re/panda_studio/internal/config"
)

func main() {
	fmt.Println("Hello!")
	if err := testS3(); err != nil {
		panic(err)
	}
}

func testS3() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	s3Config, err := config.GetConfig().GetS3Config()
	if err != nil {
		return err
	}
	bucketConfig, err := config.GetConfig().GetS3BucketsConfig()
	if err != nil {
		return err
	}

	fmt.Printf("S3 Config - %+v\n", s3Config)
	fmt.Printf("Buckets Config - %+v\n", bucketConfig)

	minioClient, err := minio.New(s3Config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, ""),
		Secure: s3Config.SslEnabled,
	})
	if err != nil {
		return err
	}

	ctx := context.TODO()

	bucketExists, err := minioClient.BucketExists(ctx, bucketConfig.ImagesBucket)
	if err != nil {
		return err
	}

	if bucketExists {
		fmt.Println("Bucket exists")
	} else {
		fmt.Println("Bucket does not exist")
	}

	return nil
}