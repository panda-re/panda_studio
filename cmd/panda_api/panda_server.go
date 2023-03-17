package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	config "github.com/panda-re/panda_studio/internal/configuration"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/api"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/middleware"
	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

type parameters struct {
	Volume   string `json:"volume"`
	Commands string `json:"commands"`
	Name     string `json:"name"`
}

type responses struct {
	Response []string `json:"response"`
}

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	if err := runServer(); err != nil {
		panic(err)
	}
}

func runServer() error {
	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	swagger, err := api.GetSwagger()
	if err != nil {
		return err
	}
	swagger.Servers = nil
	server, err := api.NewPandaStudioServer()
	if err != nil {
		return err
	}

	// r.Use(oapimiddleware.OapiRequestValidator(swagger))
	api.RegisterHandlersWithOptions(r, server, api.GinServerOptions{
		BaseURL: "/api",
	})
	//images.ImagesRegister(apiGroup.Group("/images"))

	r.POST("/panda", postRecording)

	if err := r.Run(); err != nil {
		return err
	}

	return nil
}

func postRecording(c *gin.Context) {
	var params parameters
	var response responses

	if err := c.BindJSON(&params); err != nil {
		return
	}

	response.Response, _ = startExecutor(params.Commands)

	c.IndentedJSON(http.StatusCreated, response)
}

func startExecutor(serialized_json string) ([]string, *controller.PandaAgentRecording) {
	ctx := context.Background()

	agent, err := controller.CreateDefaultDockerPandaAgent(ctx, "/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
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
	_, err = agent.StartAgent(ctx)
	if err != nil {
		panic(err)
	}

	var result []string

	var recording *controller.PandaAgentRecording
	for _, interactions := range programs {
		fmt.Printf(" %s\n", interactions)
		instructionList := interactions.Instructions
		instructions, err := parseProgram(instructionList)
		if err != nil {
			panic(err)
		}

		for _, cmd := range instructions {
			// Check Type of command and then execute backend as needed for that command.
			if cmd != nil {
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

	_, err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

	return result, recording
}

func parseProgram(instructionList string) (models.InteractionProgramInstructionList, error) {
	var interactionProgramInstructionList models.InteractionProgramInstructionList
	commandArray := strings.Split(instructionList, "\n")
	for _, cmd := range commandArray {
		res, err := parseInstruction(cmd)
		if err != nil {
			return nil, err
		}
		if res != nil {
			interactionProgramInstructionList = append(interactionProgramInstructionList, res)
		}
	}
	return interactionProgramInstructionList, nil
}

func parseInstruction(cmd string) (models.InteractionProgramInstruction, error) {
	if cmd != "" {
		cmd := strings.TrimSpace(cmd)
		if strings.HasPrefix(cmd, "#") {
			return nil, nil
		}
		instArray := strings.SplitN(cmd, " ", 2)
		switch instArray[0] {
		case "START_RECORDING":
			return &models.StartRecordingInstruction{RecordingName: instArray[1]}, nil
		case "STOP_RECORDING":
			return &models.StopRecordingInstruction{}, nil
		case "CMD":
			return &models.RunCommandInstruction{Command: instArray[1]}, nil
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
