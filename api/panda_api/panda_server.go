package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	controller "github.com/panda-re/panda_studio/internal/panda_controller"
	"net/http"
	"time"
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
	router := gin.Default()
	router.Use(cors.New(cors.Config{
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
	router.POST("/panda", postRecording)

	router.Run("localhost:8080")
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
