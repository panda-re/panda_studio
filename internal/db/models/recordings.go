package models

import (
	"github.com/panda-re/panda_studio/internal/db"
)

type Recording struct {
	ID                 db.ObjectID `bson:"_id" json:"id"`
	Name               string      `bson:"name" json:"name"`
	Description        string      `bson:"description" json:"description"`
	ImageID            db.ObjectID `bson:"recordingImage" json:"recordingImage"`
	RecordingSnapshot  Snapshot    `bson:"recordingSnapshot" json:"recordingSnapshot"`
	RecordingNondetLog NondetLog   `bson:"recordingNondetLog" json:"recordingNondetLog"`
	// Interaction list
}

type Snapshot struct {
	ID   db.ObjectID `bson:"_id" json:"_id"`
	Name string      `bson:"name" json:"name"`
}

type NondetLog struct {
	ID   db.ObjectID `bson:"id" json:"id"`
	Name string      `bson:"name" json:"name"`
}

type CreateRecordingRequest struct {
	Name          string
	ImageID       db.ObjectID
	InteractionID db.ObjectID
}
