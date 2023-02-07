package models

import (
	"encoding/json"
	"errors"

	"github.com/panda-re/panda_studio/internal/db"
	"go.mongodb.org/mongo-driver/bson"
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

func (*RunCommandInstruction) GetInstructionType() string {
	return "command"
}

type StartRecordingInstruction struct {
	RecordingName string `bson:"recording_name" json:"recording_name"`
}

func (*StartRecordingInstruction) GetInstructionType() string {
	return "start_recording"
}

type MarshalFunc func(interface{}) ([]byte, error)
type UnmarshalFunc func([]byte, interface{}) error

type discriminatedInstruction struct {
	Type string `bson:"type" json:"type"`
	InteractionProgramInstruction `bson:"args" json:"args"`
}

func (ip *InteractionProgramInstructionList) marshal(Marshal MarshalFunc) ([]byte, error) {
	typedInstructions := make([]interface{}, len(*ip))
	for i, inst := range *ip {
		typedInstructions[i] = discriminatedInstruction{
			Type: inst.GetInstructionType(),
			InteractionProgramInstruction: inst,
		}
	}

	return Marshal(typedInstructions)
}

func (ip *InteractionProgramInstructionList) unmarshal(data []byte, Unmarshal UnmarshalFunc) error {
	var types []struct { Type string `json:"type" bson:"type"` }
	err := Unmarshal(data, &types)
	if err != nil {
		return err
	}

	// todo: this needs to be generic to json and bson
	var rawMessages []json.RawMessage
	err = Unmarshal(data, &rawMessages)
	if err != nil {
		return err
	}

	instructions := make(InteractionProgramInstructionList, len(types))

	for i, msg := range rawMessages {
		item := &instructions[i]
		switch types[i].Type {
		case "command":
			*item = &RunCommandInstruction{}
		case "start_recording":
			*item = &StartRecordingInstruction{}
		default:
			return errors.New("invalid type")
		}
		
		err = Unmarshal(msg, item)
		if err != nil {
			return err
		}

	}

	*ip = instructions
	return nil
}

func (ip *InteractionProgramInstructionList) MarshalJSON() ([]byte, error) {
	return ip.marshal(json.Marshal)
}

func (ip *InteractionProgramInstructionList) MarshalBSON() ([]byte, error) {
	return ip.marshal(bson.Marshal)
}

func (ip *InteractionProgramInstructionList) UnmarshalJSON(data []byte) error {
	return ip.unmarshal(data, json.Unmarshal)
}

func (ip *InteractionProgramInstructionList) UnmarshalBSON(data []byte) (error) {
	return ip.unmarshal(data, bson.Unmarshal)
}