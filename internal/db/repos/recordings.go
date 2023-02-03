package repos

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const RECORDINGS_TABLE string = "recordings"

type RecordingRepository interface {
	ReadRecording(ctx context.Context, recordingId db.ObjectID) (*models.Recording, error)
	DeleteRecording(ctx context.Context, imageId db.ObjectID) (*models.Recording, error)
	FindRecording(ctx context.Context, id db.ObjectID) (*models.Recording, error)
	CreateRecording(ctx context.Context, obj *models.Recording) (*models.Recording, error)
}

type mongoS3RecordingRepository struct {
	coll             *mongo.Collection
	s3Client         *minio.Client
	recordingsBucket string
}

func GetRecordingRepository(ctx context.Context) (RecordingRepository, error) {
	mongoClient, err := db.GetMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}

	s3Client, err := db.GetS3Client(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoS3RecordingRepository{
		coll:             mongoClient.Collection(IMAGES_TABLE),
		s3Client:         s3Client,
		recordingsBucket: configuration.GetConfig().S3.Buckets.RecordingsBucket,
	}, nil
}

func (m mongoS3RecordingRepository) CreateRecording(ctx context.Context, obj *models.Recording) (*models.Recording, error) {
	obj.ID = db.NewObjectID()

	// insert into mongo
	result, err := m.coll.InsertOne(ctx, obj)
	if err != nil {
		return nil, err
	}

	insertedId := result.InsertedID.(primitive.ObjectID)
	obj.ID = &insertedId

	return obj, nil
}

func (m mongoS3RecordingRepository) FindRecording(ctx context.Context, id db.ObjectID) (*models.Recording, error) {
	var result models.Recording

	err := m.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &result, nil
}

func (m mongoS3RecordingRepository) ReadRecording(ctx context.Context, recordingId db.ObjectID) (*models.Recording, error) {
	recording, err := m.FindRecording(ctx, recordingId)
	if err != nil {
		return nil, err
	}

	return recording, nil
}

func (m mongoS3RecordingRepository) DeleteRecording(ctx context.Context, recordingId db.ObjectID) (*models.Recording, error) {
	recording, err := m.FindRecording(ctx, recordingId)
	if err != nil {
		return nil, err
	}

	_, err = m.coll.DeleteOne(ctx, bson.M{
		"_id": recordingId,
	})
	if err != nil {
		return nil, err
	}

	return recording, nil
}
