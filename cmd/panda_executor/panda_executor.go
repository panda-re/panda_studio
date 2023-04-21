package main

import (
	"context"
	"fmt"
	"io"
	"os"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

const DEFAULT_QCOW_SIZE = 17711104

func main() {
	// Default agent
	ctx := context.Background()
	agent, err := controller.CreateDockerPandaAgent2(ctx)
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	commands := []string{
		"uname -a",
		"ls /",
		"touch /NEW_FILE.txt",
		"ls /",
	}

	err = agent.Connect(ctx)
	if err != nil {
		panic(err)
	}

	fileReader, err := os.Open("/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		panic(err)
	}
	fileInfo, err := fileReader.Stat()
	if err != nil {
		panic(err)
	}
	err = agent.CopyFileToContainer(ctx, fileReader, fileInfo.Size(), "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		panic(err)
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
		cmdResult, err := agent.RunCommand(ctx, cmd)
		fmt.Printf("> %s\n", cmd)
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

	ndl, err := recording.OpenNdlog(ctx)
	if err != nil {
		panic(err)
	}
	dest := fmt.Sprintf("/tmp/panda-studio/%s", recording.NdlogFilename())
	out, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	nBytes, err := io.Copy(out, ndl)
	if err != nil {
		panic(err)
	}
	if nBytes == 0 {
		panic("Bad copy")
	}
	ndl.Close()

	snp, err := recording.OpenSnapshot(ctx)
	if err != nil {
		panic(err)
	}
	dest = fmt.Sprintf("/tmp/panda-studio/%s", recording.SnapshotFilename())
	out, err = os.Create(dest)
	if err != nil {
		panic(err)
	}
	nBytes, err = io.Copy(out, snp)
	if err != nil {
		panic(err)
	}
	if nBytes == 0 {
		panic("Bad copy")
	}
	snp.Close()

	fmt.Printf("Snapshot file: %s\n", recording.SnapshotFilename())
	fmt.Printf("Nondet log file: %s\n", recording.NdlogFilename())

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

	// Replay agent
	// replay_agent, err := controller.CreateDockerPandaAgent2(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer replay_agent.Close()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Starting replay")
	// replay, err := replay_agent.StartReplay(ctx, "test")
	// if err != nil {
	// 	panic(err)
	// }
	// println(replay.Serial)
	// println(replay.Replay)

	// err = replay_agent.StopAgent(ctx)
	// if err != nil {
	// 	panic(err)
	// }
}
