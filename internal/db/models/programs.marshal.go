package models

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// I hate the below code but I spent waaay too much time on it to devise
// a better solution...

type RawElement []byte

func (r RawElement) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

func (r *RawElement) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(r).UnmarshalJSON(data)
}

type MarshalFunc func(interface{}) ([]byte, error)
type UnmarshalFunc func([]byte, interface{}) error

func (ip *InteractionProgramInstructionList) marshal(Marshal MarshalFunc) ([]byte, error) {
	typedInstructions := make([]interface{}, len(*ip))
	for i, inst := range *ip {
		switch inst := inst.(type) {
		case *RunCommandInstruction:
			typedInstructions[i] = struct {
				Type string `bson:"type" json:"type"`
				*RunCommandInstruction `bson:",inline" json:",inline"`
			}{
				Type: inst.GetInstructionType(),
				RunCommandInstruction: inst,
			}
		case *StartRecordingInstruction:
			typedInstructions[i] = struct {
				Type string `bson:"type" json:"type"`
				*StartRecordingInstruction `bson:",inline" json:",inline"`
			}{
				Type: inst.GetInstructionType(),
				StartRecordingInstruction: inst,
			}
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

	var rawMessages []RawElement
	err = Unmarshal(data, &rawMessages)
	if err != nil {
		return err
	}

	instructions := make(InteractionProgramInstructionList, len(types))

	for i, msg := range rawMessages {
		item := &instructions[i]
		switch types[i].Type {
		case RunCommandInstruction{}.GetInstructionType():
			*item = &RunCommandInstruction{}
		case StartRecordingInstruction{}.GetInstructionType():
			*item = &StartRecordingInstruction{}
		case StopRecordingInstruction{}.GetInstructionType():
			*item = &StopRecordingInstruction{}
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