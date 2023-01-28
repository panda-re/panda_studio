package repos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const IMAGES_TABLE string = "images"

type ImageRepository interface {
	FindAll(ctx context.Context) ([]models.Image, error)
	FindOne(ctx context.Context, id db.ObjectID) (*models.Image, error)
	FindOneImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (*models.ImageFile, error)
	CreateImageFile(ctx context.Context, request *models.ImageFileCreateRequest) (*models.ImageFile, error)
	UploadImageFile(ctx context.Context, req *models.ImageFileUploadRequest, reader io.Reader) (*models.ImageFile, error)
	OpenImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (io.ReadCloser, error)
	DeleteImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (*models.ImageFile, error)
}

type mongoS3ImageRepository struct {
	coll *mongo.Collection
	s3Client *minio.Client
	imagesBucket string
}

func GetImageRepository(ctx context.Context) (ImageRepository, error) {
	mongoClient, err := db.GetMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}

	s3Client, err := db.GetS3Client(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoS3ImageRepository {
		coll: mongoClient.Collection(IMAGES_TABLE),
		s3Client: s3Client,
		imagesBucket: configuration.GetConfig().S3.Buckets.ImagesBucket,
	}, nil
}

func (r *mongoS3ImageRepository) FindAll(ctx context.Context) ([]models.Image, error) {
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	var images []models.Image
	if err = cursor.All(ctx, &images); err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return images, nil
}

func (r *mongoS3ImageRepository) FindOne(ctx context.Context, id db.ObjectID) (*models.Image, error) {
	var result models.Image

	err := r.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &result, nil
}

func (r *mongoS3ImageRepository) CreateImageFile(ctx context.Context, req *models.ImageFileCreateRequest) (*models.ImageFile, error) {
	newFile := models.ImageFile{
		ID: db.NewObjectID(),
		FileName: req.FileName,
		FileType: req.FileType,
		IsUploaded: false,
		Size: -1,
		Sha256: "",
	}

	_, err := r.coll.UpdateByID(ctx, req.ImageID, bson.D{
		{"$push", bson.D{
			{"files", newFile},
		},
	}})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &newFile, err
}

func (r *mongoS3ImageRepository) FindOneImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (*models.ImageFile, error) {
	var img models.Image
	err := r.coll.FindOne(ctx, bson.M{
		"_id": imageId,
		"files": bson.D{{"$elemMatch",
			bson.D{{"_id", fileId}},
		}},
	}, options.FindOne().SetProjection(bson.D{
		// Filters the files to just the one we want
		{"files.$", 1},
	})).Decode(&img)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	// Our query should only return one element but just in case
	if len(img.Files) > 1 {
		return nil, errors.New("Something is off with the query")
	}

	imgFile := img.Files[0]
	imgFile.ImageID = imageId
	
	return imgFile, nil
}

func (r *mongoS3ImageRepository) getObjectName(imageId db.ObjectID, file *models.ImageFile) string {
	objectName := fmt.Sprintf("%s/%s", imageId.Hex(), file.FileName)
	return objectName
}

func (r *mongoS3ImageRepository) UploadImageFile(ctx context.Context, req *models.ImageFileUploadRequest, reader io.Reader) (*models.ImageFile, error) {
	imgFile, err := r.FindOneImageFile(ctx, req.ImageId, req.FileId)
	if err != nil {
		return nil, err
	}

	objectName := r.getObjectName(req.ImageId, imgFile)

	hasher := sha256.New()
	hashReader := util.NewReaderHasher(reader, hasher)

	// Upload file and compute hash
	obj, err := r.s3Client.PutObject(ctx, r.imagesBucket, objectName, hashReader, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "file upload failed")
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Update the file information
	_, err = r.coll.UpdateOne(ctx, bson.M{
		"_id": req.ImageId,
		"files._id": req.FileId,
	}, bson.D{
		{"$set", bson.M{
			"files.$.is_uploaded": true,
			"files.$.size": obj.Size,
			"files.$.sha256": hash,
		}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	imgFile, err = r.FindOneImageFile(ctx, req.ImageId, req.FileId)
	if err != nil {
		return nil, err
	}

	return imgFile, nil
}


func (r *mongoS3ImageRepository) OpenImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (io.ReadCloser, error) {
	imgFile, err := r.FindOneImageFile(ctx, imageId, fileId)
	if err != nil {
		return nil, err
	}

	objectName := r.getObjectName(imageId, imgFile)

	obj, err := r.s3Client.GetObject(ctx, r.imagesBucket, objectName, minio.GetObjectOptions{
		Checksum: true,
	})
	if err != nil {
		return nil, err
	}

	return obj, nil
}


func (r *mongoS3ImageRepository) DeleteImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (*models.ImageFile, error) {
	imgFile, err := r.FindOneImageFile(ctx, imageId, fileId)
	if err != nil {
		return nil, err
	}

	// delete file
	objName := r.getObjectName(imageId, imgFile)

	err = r.s3Client.RemoveObject(ctx, r.imagesBucket, objName, minio.RemoveObjectOptions{})
	if err != nil {
		return nil, err
	}

	// delete from db
	_, err = r.coll.UpdateOne(ctx, bson.M{
		"_id": imageId,
		"files._id": fileId,
	}, bson.D{
		{"$pull", bson.M{
			"files": bson.M{
				"_id": fileId,
			},
		}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return imgFile, nil
}
