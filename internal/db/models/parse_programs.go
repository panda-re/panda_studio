package models

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func ParseInteractionProgram(instructionList string) (InteractionProgramInstructionList, error) {
	var interactionProgramInstructionList InteractionProgramInstructionList
	commandArray := strings.Split(instructionList, "\n")
	for _, cmd := range commandArray {
		res, err := ParseInteractionProgramInstruction(cmd)
		if err != nil {
			return nil, err
		}
		if res != nil {
			interactionProgramInstructionList = append(interactionProgramInstructionList, res)
		}
	}
	return interactionProgramInstructionList, nil
}

func ParseInteractionProgramInstruction(cmd string) (InteractionProgramInstruction, error) {
	if cmd != "" {
		cmd := strings.TrimSpace(cmd)
		if strings.HasPrefix(cmd, "#") {
			return nil, nil
		}
		instArray := strings.SplitN(cmd, " ", 2)
		switch strings.ToLower(instArray[0]) {
		case "start_recording":
			return &StartRecordingInstruction{RecordingName: instArray[1]}, nil
		case "stop_recording":
			return &StopRecordingInstruction{}, nil
		case "command":
			return &RunCommandInstruction{Command: instArray[1]}, nil
		case "filesystem":
			fmt.Printf("Filesystem placeholder\n")
			return nil, errors.New("Filesystem interactions not yet supported")
		case "network":
			fmt.Printf("Network placeholder\n")
			return nil, errors.New("Network interactions not yet supported")
		default:
			fmt.Printf("Incorrect Command Type %s, Correct options can be found in the commands.md file\n", instArray[0])
			return nil, errors.Errorf("Incorrect Command Type %s, Correct options can be found in the commands.md file\n", instArray[0])
		}
	}

	return nil, nil
}