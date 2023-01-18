package images

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Files []ImageFile `bson:"files" json:"files"`
	Config *ImageConfiguration `bson:"config" json:"config"`
}

type ImageFile struct {

}

type ImageConfiguration struct {}