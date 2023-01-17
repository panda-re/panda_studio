package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const IMAGES_TABLE string = "images"

type Image struct {
	ID ObjectID `bson:"_id"`
	Name string `bson:"name"`
	Description string `bson:"description"`
	Files []ImageFile `bson:"files"`
	Config *ImageConfiguration `bson:"config"`
}

type ImageFile struct {

}

type ImageConfiguration struct {}

type ImageRepository interface {
	//FindImage()
}

type mongoS3ImageRespository struct {
	coll *mongo.Collection
}

func GetRepository(ctx context.Context) (ImageRepository, error) {
	db, err := GetMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoS3ImageRespository {
		coll: db.Collection(IMAGES_TABLE),
	}, nil
}
