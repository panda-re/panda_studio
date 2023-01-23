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

	for _, cmd := range commands {
		fmt.Printf("> %s\n", cmd)
		cmdResult, err := agent.RunCommand(ctx, cmd)
		if err != nil {
			fmt.Println("Test Passed. Prevented running Commands before Panda Starts")
			break
		} else {
			fmt.Println("Test Failed. Command Run Without Starting PANDA")
			panic(err)
		}
		fmt.Printf("%s\n", cmdResult.Logs)
	}

	err = agent.StopAgent(ctx)
	if err != nil {
		fmt.Println("Test Passed. Prevented Stopping PANDA before PANDA Starts")
	} else {
		fmt.Println("Test Failed. Stopped Nonexistent PANDA")
		panic(err)
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

	fmt.Println("Stopping recording before starting it")
	recording, err := agent.StopRecording(ctx)
	if err != nil {
		fmt.Println("Test Passed. Prevented Stopping of a non-existent recording")
	} else {
		fmt.Println("Test Failed. Stopped Nothing From ever Recording")
		panic(err)
	}

	fmt.Println("Starting recording with Spaces as the name, it will work")
	if err := agent.StartRecording(ctx, "         "); err != nil {
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
	recording, err = agent.StopRecording(ctx)
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
