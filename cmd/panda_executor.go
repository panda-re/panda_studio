package main

import (
	"context"
	"fmt"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

func main() {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	//defer cancel()

	// agent, err := controller.CreateDefaultGrpcPandaAgent()
	agent, err := controller.CreateDefaultDockerPandaAgent(ctx)
	if err != nil {
		panic(err)
	}
	// defer agent.Close()

	if err != nil {
		panic(err)
	}

	commands := []string{
		"uname -a",
		"ls /",
		"touch /NEW_FILE.txt",
		"ls /",
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
