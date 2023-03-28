package models

import (
	"github.com/panda-re/panda_studio/internal/db"
)

type InteractionProgramInstructionList []InteractionProgramInstruction

type InteractionProgram struct {
	ID           db.ObjectID `bson:"_id" json:"id"`
	Name         string      `bson:"name" json:"name"`
	Instructions string      `bson:"instructions" json:"instructions"`
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

type FilesystemInstruction struct {
}

func (FilesystemInstruction) GetInstructionType() string {
	return "filesystem"
}

type NetworkInstruction struct {
	SocketType string `bson:"sock_type" json:"sock_type"`
	Port       int    `bson:"port" json:"port"`
	PacketType string `bson:"packet_type" json:"packet_type"`
	PacketData string `bson:"packet_data" json:"packet_data"`
}

func (NetworkInstruction) GetInstructionType() string {
	return "network"
}
