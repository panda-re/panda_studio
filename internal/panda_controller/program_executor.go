package panda_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/db/repos"
	"github.com/pkg/errors"
)

type PandaProgramExecutor struct {
	imageRepo repos.ImageRepository
}

type PandaProgramExecutorJob struct {
	agent PandaAgent
	image *models.Image
	program *models.InteractionProgram
}

type PandaProgramExecutorOptions struct {
	Image *models.Image
	Program *models.InteractionProgram
	Instructions models.InteractionProgramInstructionList
	ImageFileReader io.Reader
}

func (p *PandaProgramExecutor) NewExecutorJob(opts *PandaProgramExecutorOptions) (*PandaProgramExecutorJob, error) {
	// Basic rundown of what will happen:
	// 1. pull information from the database
	//    - making this the caller's responsibility
	//	  - passed in as 'opts'
	// 2. open a stream to the file in blob storage
	//    - caller's responsibility for now
	// 3. create a panda instance using that file
	// 4. start the agent
	// 5. send the commands to the agent
	//    - offer an interface for real-time feedback, even if we don't currently use it
	//    - keep track of any recording files that are created
	// 6. stop the agent
	// 7. upload the recording files to blob storage
	return nil, errors.New("not implemented")
}

func startExecutor(serialized_json string) ([]string, *PandaAgentRecording) {
	ctx := context.Background()

	// todo: change this method to take in a `Reader` interface instead of a path
	agent, err := CreateDefaultDockerPandaAgent(ctx, "/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	if err != nil {
		panic(err)
	}

	var programs []models.InteractionProgram

	err = json.Unmarshal([]byte(serialized_json), &programs)
	if err != nil {
		panic(err)
	}

	// Start Agent assuming that we are not running a replay
	fmt.Println("Starting agent")
	err = agent.StartAgent(ctx)
	if err != nil {
		panic(err)
	}

	var result []string

	var recording *PandaAgentRecording
	for _, interactions := range programs {
		fmt.Printf(" %s\n", interactions)
		instructionList := interactions.Instructions
		instructions, err := models.ParseInteractionProgram(instructionList)
		if err != nil {
			panic(err)
		}

		for _, cmd := range instructions {
			// Check Type of command and then execute backend as needed for that command.
			if cmd != nil {
				// todo: I think we should make this polymorphic
				switch cmd.GetInstructionType() {
				case "start_recording":
					// Since we have a start recording command, we have to type cast cmd to a pointer for a StartRecordingInstruction from the models package
					err := agent.StartRecording(ctx, cmd.(*models.StartRecordingInstruction).RecordingName)
					if err != nil {
						panic(err)
					}
					break
				case "stop_recording":
					recording, err = agent.StopRecording(ctx)
					if err != nil {
						panic(err)
					}
					break
				case "command":
					cmdResult, err := agent.RunCommand(ctx, cmd.(*models.RunCommandInstruction).Command)
					if err != nil {
						panic(err)
					}
					fmt.Printf(" %s\n", cmdResult)
					result = append(result, cmdResult.Logs+"\n")
					break
				case "filesystem":
					fmt.Printf("Filesystem placeholder\n")
					break
				case "network":
					fmt.Printf("Network Placeholder\n")
					fmt.Printf("%s\n", cmd.(*models.NetworkInstruction).SocketType)
					fmt.Printf("%d\n", cmd.(*models.NetworkInstruction).Port)
					fmt.Printf("%s\n", cmd.(*models.NetworkInstruction).PacketType)
					fmt.Printf("%s\n", cmd.(*models.NetworkInstruction).PacketData)
					break
				default:
					fmt.Printf("Incorrect Command Type, Correct options can be found in the commands.md file")
					break
				}
			}
		}

	}

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

	return result, recording
}
