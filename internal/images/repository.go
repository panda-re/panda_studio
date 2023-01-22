package images

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/panda-re/panda_studio/internal/db"
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
	CreateImageFile(ctx context.Context, request *ImageFileCreateRequest) (*ImageFile, error)
	UploadImageFile(ctx context.Context, req *ImageFileUploadRequest, reader io.Reader) (*ImageFile, error)
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
	var result Image

	err := r.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &result, nil
}

func (r *mongoS3ImageRespository) CreateImageFile(ctx context.Context, req *ImageFileCreateRequest) (*ImageFile, error) {
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

func (r *mongoS3ImageRespository) UploadImageFile(ctx context.Context, req *ImageFileUploadRequest, reader io.Reader) (*ImageFile, error) {
	var img Image
	err := r.coll.FindOne(ctx, bson.D{
		{"_id", req.ImageId},
		{"files._id", req.FileId},
	}, options.FindOne().SetProjection(bson.D{
		{"files", bson.E{"$elemMatch", bson.E{"_id", req.FileId}}},
	})).Decode(&img)

	if err != nil {
		return nil, err
	}

	var imgFile *ImageFile
	for _, imgF := range img.Files {
		if imgF == nil || req.FileId == nil {
			continue
		}
		if *imgF.ID == *req.FileId {
			imgFile = imgF
		}
	}
	if imgFile == nil {
		return nil, errors.New("Something is wrong")
	}

	objectName := fmt.Sprintf("%s/%s", req.ImageId.Hex(), imgFile.FileName)

	hasher := sha256.New()
	readPipe, writePipe := io.Pipe()
	multiWriter := io.MultiWriter(writePipe, hasher)

	go func() {
		io.Copy(multiWriter, reader)
		readPipe.Close()
	}()

	// Upload file and compute hash
	// todo: don't hardcode bucket name!
	info, err := r.s3Client.PutObject(ctx, "images", objectName, readPipe, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "file upload failed")
	}
	fmt.Printf("%+v", info)

	return nil, errors.New("Not implemented")
}
