package main

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	config "github.com/panda-re/panda_studio/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("Hello!")
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	if err := testMongo(); err != nil {
		panic(err)
	}
	if err := testS3(); err != nil {
		panic(err)
	}
}

func testS3() error {

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
		fmt.Println("Bucket does not exist... creating")
		if err := minioClient.MakeBucket(ctx, bucketConfig.ImagesBucket, minio.MakeBucketOptions{ObjectLocking: true}); err != nil {
			return err
		}
	}

	for obj := range minioClient.ListObjects(ctx, bucketConfig.ImagesBucket, minio.ListObjectsOptions{
		Prefix: "08192a3b4c5d6e7f/",
	}) {
		if obj.Err != nil {
			return obj.Err
		}
		fmt.Printf("Bucket Objects - %+v\n", obj)
	}


	return nil
}

func testMongo() error {
	ctx := context.TODO()

	mongoUri, err := config.GetConfig().GetMongoConfig()
	if err != nil {
		return err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri.Uri))
	if err != nil {
		return err
	}

	fmt.Println("Pinging Mongodb!")
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	userCollection := client.Database("testing").Collection("users")
	user := bson.D{{Key: "fullName", Value: "User 1"}, {Key: "age", Value: 30}}
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	fmt.Printf("Mongo Insert Result: %+v\n", result)
	return nil

	result2, err := userCollection.DeleteOne(ctx, bson.D{{"_id", result.InsertedID}})
	if err != nil {
		return err
	}

	fmt.Printf("Mongo Delete Result: %+v\n", result2)
	return nil
}