package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/api"
	config "github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/middleware"
	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

type command struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Name    string `json:"name"`
}

type parameters struct {
	Volume   string    `json:"volume"`
	Commands []command `json:"commands"`
	Name     string    `json:"name"`
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

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

func startExecutor(commands []command) ([]string, string) {
	ctx := context.Background()

	// Create the docker contaier
	agent, err := controller.CreateDefaultDockerPandaAgent(ctx, "")
	if err != nil {
		panic(err)
	}
	defer agent.Close()

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

	var recording string
	//var cmdResult string

	for _, cmd := range commands {
		fmt.Printf("> %s\n", cmd)
		// Check Type of command and then execute backend as needed for that command.
		switch cmd.Type {
		case "Recording":
			if cmd.Command == "start" {
				cmdResult, err := agent.StartRecording(ctx, cmd.Name)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", cmdResult.Logs)
				result = append(result, cmdResult.Logs+"\n")
			} else if cmd.Command == "stop" {
				recording, err = agent.StopRecording(ctx)
				if err != nil {
					panic(err)
				}
			} else {
				panic("Error, Recording Instruction Type with incorrect instruction")
			}
			break
		case "Serial":
			cmdResult, err := agent.RunCommand(ctx, cmd.Command)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", cmdResult.Logs)
			result = append(result, cmdResult.Logs+"\n")
			break
		case "Filesystem":
			break
		case "Network":
			break
		default:
			panic("Incorrect Command Type, Correct Options are: Recording, Serial, Filesystem or Network")
		}

	}

	return result, recording
}
