package images

import (
	"context"

	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const IMAGES_TABLE string = "images"

type ImageRepository interface {
	FindAll(ctx context.Context) ([]Image, error)
}

type mongoS3ImageRespository struct {
	coll *mongo.Collection
}

func GetRepository(ctx context.Context) (ImageRepository, error) {
	db, err := db.GetMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoS3ImageRespository {
		coll: db.Collection(IMAGES_TABLE),
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
