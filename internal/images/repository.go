package images

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const IMAGES_TABLE string = "images"

type Repository [T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindOne(ctx context.Context, id db.ObjectID) (*T, error)
}

type ImageRepository interface {
	Repository[Image]
}

type mongoS3ImageRespository struct {
	coll *mongo.Collection
	s3Client *minio.Client
}

func GetRepository(ctx context.Context) (ImageRepository, error) {
	mongoClient, err := db.GetMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}

	s3Client, err := db.GetS3Client(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoS3ImageRespository {
		coll: mongoClient.Collection(IMAGES_TABLE),
		s3Client: s3Client,
	}, nil
}

func (r *mongoS3ImageRespository) FindAll(ctx context.Context) ([]Image, error) {
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	var images []Image
	if err = cursor.All(ctx, &images); err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return images, nil
}

func (r *mongoS3ImageRespository) FindOne(ctx context.Context, id db.ObjectID) (*Image, error) {
	imageId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}
	var result Image

	err = r.coll.FindOne(ctx, bson.D{{"_id", imageId}}).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &result, nil
}
