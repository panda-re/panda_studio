package images

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const IMAGES_TABLE string = "images"

type Repository [T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindOne(ctx context.Context, id db.ObjectID) (*T, error)
}

type ImageRepository interface {
	Repository[Image]
	FindOneImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (*ImageFile, error)
	CreateImageFile(ctx context.Context, request *ImageFileCreateRequest) (*ImageFile, error)
	UploadImageFile(ctx context.Context, req *ImageFileUploadRequest, reader io.Reader) (*ImageFile, error)
	OpenImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (io.ReadSeekCloser, error)
}

type mongoS3ImageRepository struct {
	coll *mongo.Collection
	s3Client *minio.Client
	imagesBucket string
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

	return &mongoS3ImageRepository {
		coll: mongoClient.Collection(IMAGES_TABLE),
		s3Client: s3Client,
		imagesBucket: configuration.GetConfig().S3.Buckets.ImagesBucket,
	}, nil
}

func (r *mongoS3ImageRepository) FindAll(ctx context.Context) ([]Image, error) {
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

func (r *mongoS3ImageRepository) FindOne(ctx context.Context, id db.ObjectID) (*Image, error) {
	var result Image

	err := r.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &result, nil
}

func (r *mongoS3ImageRepository) CreateImageFile(ctx context.Context, req *ImageFileCreateRequest) (*ImageFile, error) {
	newFile := ImageFile{
		ID: db.NewObjectID(),
		FileName: req.FileName,
		FileType: req.FileType,
		IsUploaded: false,
		SHA256: "",
	}

	_, err := r.coll.UpdateByID(ctx, req.ImageID, bson.D{{"$push", bson.D{
		{"files", newFile},
	}}})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &newFile, err
}

func (r *mongoS3ImageRepository) FindOneImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (*ImageFile, error) {
	var img Image
	err := r.coll.FindOne(ctx, bson.D{
		{"_id", imageId},
		{"files",
			bson.D{{"$elemMatch",
				bson.D{{"_id", fileId}},
			}},
		},
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
	return img.Files[0], nil
}

func (r *mongoS3ImageRepository) getObjectName(imageId db.ObjectID, file *ImageFile) string {
	objectName := fmt.Sprintf("%s/%s", imageId.Hex(), file.FileName)
	return objectName
}

func (r *mongoS3ImageRepository) UploadImageFile(ctx context.Context, req *ImageFileUploadRequest, reader io.Reader) (*ImageFile, error) {
	imgFile, err := r.FindOneImageFile(ctx, req.ImageId, req.FileId)
	if err != nil {
		return nil, err
	}

	objectName := r.getObjectName(req.ImageId, imgFile)

	hasher := sha256.New()
	hashReader := util.NewReaderHasher(reader, hasher)

	// Upload file and compute hash
	_, err = r.s3Client.PutObject(ctx, r.imagesBucket, objectName, hashReader, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "file upload failed")
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Update the file information
	_, err = r.coll.UpdateOne(ctx, bson.D{
		{"_id", req.ImageId},
		{"files._id", req.FileId},
	}, bson.D{
		{"$set", bson.D{
			bson.E{"files.$.is_uploaded", true},
			bson.E{"files.$.sha256", hash},
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


func (r *mongoS3ImageRepository) DownloadImageFile(ctx context.Context, imageId db.ObjectID, fileId db.ObjectID) (io.ReadSeekCloser, error) {
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
