package repos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
)

const RECORDINGS_TABLE string = "recordings"

type RecordingRepository interface {
	ReadRecording(ctx context.Context, recordingId db.ObjectID) (*models.Recording, error)
	DeleteRecording(ctx context.Context, recordingId db.ObjectID) (*models.Recording, error)
	FindRecording(ctx context.Context, id db.ObjectID) (*models.Recording, error)
	FindAllRecordings(ctx context.Context) ([]models.Recording, error)
	CreateRecording(ctx context.Context, obj *models.Recording) (*models.Recording, error)
	FindRecordingFile(ctx context.Context, recordingId db.ObjectID, fileId db.ObjectID) (*models.RecordingFile, error)
	DeleteRecordingFile(ctx context.Context, recordingId db.ObjectID, fileId db.ObjectID) (*models.RecordingFile, error)
	CreateRecordingFile(ctx context.Context, req *models.CreateRecordingFileRequest) (*models.RecordingFile, error)
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
		coll:             mongoClient.Collection(RECORDINGS_TABLE),
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

func (m *mongoS3RecordingRepository) FindAllRecordings(ctx context.Context) ([]models.Recording, error) {
	cursor, err := m.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	var recordings []models.Recording
	if err = cursor.All(ctx, &recordings); err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return recordings, nil
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

	for _, recordingFile := range recording.Files {
		_, err := m.DeleteRecordingFile(ctx, recordingId, recordingFile.ID)
		if err != nil {
			return nil, err
		}
	}

	_, err = m.coll.DeleteOne(ctx, bson.M{
		"_id": recordingId,
	})
	if err != nil {
		return nil, err
	}

	return recording, nil
}

func (m *mongoS3RecordingRepository) CreateRecordingFile(ctx context.Context, req *models.CreateRecordingFileRequest) (*models.RecordingFile, error) {
	newRecordingFile := models.RecordingFile{
		ID:         db.NewObjectID(),
		Name:       req.Name,
		FileType:   req.FileType,
		IsUploaded: false,
		Size:       -1,
		Sha256:     "",
	}

	_, err := m.coll.UpdateByID(ctx, req.RecordingID, bson.D{
		{"$push", bson.D{
			{"files", newRecordingFile},
		},
		}})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return &newRecordingFile, err
}

func (m *mongoS3RecordingRepository) UploadRecordingFile(ctx context.Context, req *models.UploadRecordingFileRequest, reader io.Reader) (*models.RecordingFile, error) {
	recordingFile, err := m.FindRecordingFile(ctx, req.RecordingID, req.FileID)
	if err != nil {
		return nil, err
	}

	objectName := m.getObjectName(req.RecordingID, recordingFile)

	hasher := sha256.New()
	hashReader := util.NewReaderHasher(reader, hasher)

	obj, err := m.s3Client.PutObject(ctx, m.recordingsBucket, objectName, hashReader, -1, minio.PutObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "file upload failed")
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	_, err = m.coll.UpdateOne(ctx, bson.M{
		"_id":       req.RecordingID,
		"files._id": req.FileID,
	}, bson.D{
		{"$set", bson.M{
			"files.$.is_uploaded": true,
			"files.$.size":        obj.Size,
			"files.$.sha256":      hash,
		}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	recordingFile, err = m.FindRecordingFile(ctx, req.RecordingID, req.FileID)
	if err != nil {
		return nil, err
	}

	return recordingFile, nil
}

func (m *mongoS3RecordingRepository) FindRecordingFile(ctx context.Context, recordingId db.ObjectID, fileId db.ObjectID) (*models.RecordingFile, error) {
	var recording models.Recording
	err := m.coll.FindOne(ctx, bson.M{
		"_id": recordingId,
		"files": bson.D{{"$elemMatch",
			bson.D{{"_id", fileId}},
		}},
	}, options.FindOne().SetProjection(bson.D{
		{"files.$", 1},
	})).Decode(&recording)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	if len(recording.Files) > 1 {
		return nil, errors.New("Something is off with the query")
	}

	recordingFile := recording.Files[0]

	return recordingFile, nil
}

func (m *mongoS3RecordingRepository) DeleteRecordingFile(ctx context.Context, recordingId db.ObjectID, fileId db.ObjectID) (*models.RecordingFile, error) {
	recordingFile, err := m.FindRecordingFile(ctx, recordingId, fileId)
	if err != nil {
		return nil, err
	}

	objName := m.getObjectName(recordingId, recordingFile)

	err = m.s3Client.RemoveObject(ctx, m.recordingsBucket, objName, minio.RemoveObjectOptions{})
	if err != nil {
		return nil, err
	}

	_, err = m.coll.UpdateOne(ctx, bson.M{
		"_id":       recordingId,
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

	return recordingFile, nil
}

func (m *mongoS3RecordingRepository) getObjectName(recordingId db.ObjectID, file *models.RecordingFile) string {
	objectName := fmt.Sprintf("%s/%s", recordingId.Hex(), file.Name)
	return objectName
}
