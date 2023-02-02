package main

import (
	"context"
	"fmt"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

func main() {
	// Default agent
	ctx := context.Background()
	agent, err := controller.CreateDefaultDockerPandaAgent(ctx, "/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
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

	fmt.Println("Starting recording")
	if err := agent.StartRecording(ctx, "test"); err != nil {
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

	// Replay agent
	replay_agent, err := controller.CreateReplayDockerPandaAgent(ctx)
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting replay")
	replay, err := replay_agent.StartReplayAgent(ctx, "test")
	if err != nil {
		panic(err)
	}
	println(replay.Serial)
	println(replay.Replay)

	err = replay_agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}
