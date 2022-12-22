package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	controller "github.com/panda-re/panda_studio/internal/panda_controller"
	"net/http"
)

type parameters struct {
	Volume   string   `json:"volume"`
	Commands []string `json:"commands"`
	Name     string   `json:"name"`
}

func main() {
	router := gin.Default()

	router.POST("/panda", postRecording)

	router.Run("localhost:8080")
}

func postRecording(c *gin.Context) {
	var params parameters

	if err := c.BindJSON(&params); err != nil {
		return
	}

	startExecutor(params.Commands)

	c.IndentedJSON(http.StatusCreated, params)

	fmt.Println(params.Commands)
}

func startExecutor(commands []string) {
	ctx := context.Background()

	agent, err := controller.CreateDefaultDockerPandaAgent(ctx)
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

	for _, cmd := range commands {
		fmt.Printf("> %s\n", cmd)
		cmdResult, err := agent.RunCommand(ctx, cmd)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", cmdResult.Logs)
	}

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}
