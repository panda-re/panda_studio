package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/api"
	config "github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/middleware"
	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

type parameters struct {
	Volume   string   `json:"volume"`
	Commands []string `json:"commands"`
	Name     string   `json:"name"`
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

	response.Response = startExecutor(params.Commands)

	c.IndentedJSON(http.StatusCreated, response)
}

func startExecutor(commands []string) []string {
	ctx := context.Background()

	agent, err := controller.CreateDefaultDockerPandaAgent(ctx, "")
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting agent")
	err = agent.StartAgent(ctx)
	if err != nil {
		panic(err)
	}

	var result []string

	for _, cmd := range commands {
		fmt.Printf("> %s\n", cmd)
		cmdResult, err := agent.RunCommand(ctx, cmd)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", cmdResult.Logs)
		result = append(result, cmdResult.Logs+"\n")
	}

	return result
}
