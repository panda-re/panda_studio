package models

import (
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
)

type InteractionProgram struct {
	ID db.ObjectID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Instructions []InteractionProgramInstruction `bson:"instructions" json:"instructions"`
}

type InteractionProgramInstruction interface {
	Discriminator() string
}

type RunCommandInstruction struct {
	Command string `bson:"command" json:"command"`
}

func (*RunCommandInstruction) Discriminator() string {
	return "command"
}

type StartRecordingInstruction struct {
	RecordingName string `bson:"recording_name" json:"recording_name"`
}

func (*StartRecordingInstruction) Discriminator() string {
	return "start_recording"
}

func (ip *InteractionProgram) MarshalJSON() ([]byte, error) {

	return nil, errors.New("Not implemented")
}

func (ip *InteractionProgram) UnmarshalJSON(b []byte) error {

	return errors.New("Not implemented")
}