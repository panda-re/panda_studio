package main

import (
	"context"
	"fmt"
	"io"

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
	stream, err := agent.StartAgent(ctx)
	if err != nil {
		panic(err)
	}
	if stream == nil {
		panic("no stream")
	}
	// Required for proper startup
	resp, err := stream.Recv()
	if err != nil {
		panic(err)
	} else if resp.Execution != "" {
		print("first resp not empty")
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

	log, err := agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Log file: %s\n", log.GetLogFileName())

	for {
		resp, err = stream.Recv()
		if err == io.EOF {
			break
		}
		// Errors out when done
		if err != nil {
			panic(err)
		}
		print(resp.Execution)
	}

	// Replay agent
	replay_agent, err := controller.CreateReplayDockerPandaAgent(ctx)
	if err != nil {
		panic(err)
	}
	defer replay_agent.Close()

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting replay")
	replay_stream, err := replay_agent.StartReplayAgent(ctx, "test")
	if err != nil {
		panic(err)
	}
	if replay_stream == nil {
		panic("no stream")
	}
	// Required for proper startup
	replay_resp, err := replay_stream.Recv()
	if err != nil {
		panic(err)
	} else if replay_resp.Replay != "" || replay_resp.Serial != "" {
		print("first resp not empty")
	}

	for {
		replay_resp, err = replay_stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if replay_resp.Replay != "" {
			print(replay_resp.Replay)
		}
		if replay_resp.Serial != "" {
			print(replay_resp.Serial)
		}
	}

	log, err = replay_agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Log file: %s\n", log.GetLogFileName())
}
