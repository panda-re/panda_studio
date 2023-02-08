package models

import (
	"github.com/panda-re/panda_studio/internal/db"
)

type Recording struct {
	ID             db.ObjectID      `bson:"_id" json:"id"`
	Name           string           `bson:"name" json:"name"`
	Description    string           `bson:"description" json:"description"`
	ImageID        db.ObjectID      `bson:"recordingImage" json:"recordingImage"`
	RecordingFiles []*RecordingFile `bson:"recordingSnapshot" json:"recordingSnapshot"`
	// Interaction list
}

type RecordingFile struct {
	ID         db.ObjectID `bson:"_id" json:"_id"`
	Name       string      `bson:"name" json:"name"`
	FileType   string      `bson:"file_type" json:"file_type"`
	IsUploaded bool        `bson:"is_uploaded" json:"is_uploaded"`
	Size       int64       `bson:"size" json:"size"`
	Sha256     string      `bson:"sha256" json:"sha256,omitempty"`
}

type CreateRecordingRequest struct {
	Name          string
	ImageID       db.ObjectID
	InteractionID db.ObjectID
}

type CreateRecordingFileRequest struct {
	Name        string
	RecordingID db.ObjectID
	FileType    string `bson:"file_type" json:"file_type"`
}

type UploadRecordingFileRequest struct {
	RecordingID db.ObjectID
	FileID      db.ObjectID
}
