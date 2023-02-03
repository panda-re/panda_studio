package models

import "github.com/panda-re/panda_studio/internal/db"

type Recording struct {
	ID             db.ObjectID `bson:"_id" json:"id"`
	Name           string      `bson:"name" json:"name"`
	Description    string      `bson:"description" json:"description"`
	RecordingImage Image       `bson:"recordingImage" json:"recordingImage"`
	// Interaction list
}
