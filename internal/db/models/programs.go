package models

import (
	"github.com/panda-re/panda_studio/internal/db"
)

type InteractionProgramInstructionList []InteractionProgramInstruction

type InteractionProgram struct {
	ID db.ObjectID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Instructions InteractionProgramInstructionList `bson:"instructions" json:"instructions"`
}

type InteractionProgramInstruction interface {
	GetInstructionType() string
}

type RunCommandInstruction struct {
	Command string `bson:"command" json:"command"`
}

func (RunCommandInstruction) GetInstructionType() string {
	return "command"
}

type StartRecordingInstruction struct {
	RecordingName string `bson:"recording_name" json:"recording_name"`
}

func (StartRecordingInstruction) GetInstructionType() string {
	return "start_recording"
}

type StopRecordingInstruction struct {

}

func (StopRecordingInstruction) GetInstructionType() string {
	return "stop_recording"
}
