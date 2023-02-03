package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	config "github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Hello!")
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config.GetConfig())

	ctx := context.TODO()

	client, err := db.GetMongoDatabase(ctx)
	if err != nil {
		panic(err)
	}

	if err := testMongo(ctx, client); err != nil {
		panic(err)
	}
	if err := testMongoGridFS(ctx, client); err != nil {
		panic(err)
	}
	if err := testS3(); err != nil {
		panic(err)
	}
}

func connectMinio() (*minio.Client, error) {
	s3Config := config.GetConfig().S3

	minioClient, err := minio.New(s3Config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, ""),
		Secure: s3Config.SslEnabled,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, err
}

func testS3() error {

	minioClient, err := connectMinio()
	if err != nil {
		return nil
	}

	bucketConfig := config.GetConfig().S3.Buckets

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

func testMongo(ctx context.Context, client *mongo.Database) error {
	userCollection := client.Collection("users")
	user := bson.D{{Key: "fullName", Value: "User 1"}, {Key: "age", Value: 30}}
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	fmt.Printf("Mongo Insert Result: %+v\n", result)

	result2, err := userCollection.DeleteOne(ctx, bson.D{{"_id", result.InsertedID}})
	if err != nil {
		return err
	}

	fmt.Printf("Mongo Delete Result: %+v\n", result2)
	return nil
}

func testMongoGridFS(ctx context.Context, db *mongo.Database) error {
	opts := options.GridFSBucket().SetName("files")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		return nil
	}

	file, err := os.Open("./random.bin")
	if err != nil {
		return err
	}
	uploadOpts := options.GridFSUpload().SetMetadata(bson.D{{"tag", "first"}})

	objectId, err := bucket.UploadFromStream("folder1/random.bin", io.Reader(file), uploadOpts)
	if err != nil {
		return err
	}

	fmt.Printf("New file uploaded with ID %s\n", objectId)

	return nil
}