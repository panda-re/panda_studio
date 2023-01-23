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
	defer agent.Close()

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
	// TODO tests
	fmt.Println("Starting Second Panda agent")
	err = agent.StartAgent(ctx)
	if err != nil {
		fmt.Println("Test Passed. Prevented second PANDA agent creation")
	} else {
		fmt.Println("Test Failed. Created second PANDA agent creation")
		panic(err)
	}

	fmt.Println("Starting recording")
	if err := agent.StartRecording(ctx, "test"); err != nil {
		panic(err)
	}
	fmt.Println("Starting second recording")
	if err := agent.StartRecording(ctx, "testing"); err != nil {
		fmt.Println("Test Passed. Prevented starting a second recording")
	} else {
		fmt.Println("Test Failed. Started a second recording")
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

	fmt.Println("Stopping recording")
	recording, err := agent.StopRecording(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Snapshot file: %s\n", recording.GetSnapshotFileName())
	fmt.Printf("Nondet log file: %s\n", recording.GetNdlogFileName())

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}
