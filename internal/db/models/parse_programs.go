package models

import (
	"errors"
	"fmt"
	"strings"
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
		switch instArray[0] {
		case "START_RECORDING":
			return &StartRecordingInstruction{RecordingName: instArray[1]}, nil
		case "STOP_RECORDING":
			return &StopRecordingInstruction{}, nil
		case "CMD":
			return &RunCommandInstruction{Command: instArray[1]}, nil
		case "filesystem":
			fmt.Printf("Filesystem placeholder\n")
			return nil, errors.New("Filesystem interactions not yet supported")
		case "network":
			fmt.Printf("Network placeholder\n")
			return nil, errors.New("Network interactions not yet supported")
		default:
			fmt.Printf("Incorrect Command Type, Correct options can be found in the commands.md file")
			return nil, errors.New("Incorrect Command Type, Correct options can be found in the commands.md file")
		}
	}

	return nil, nil
}