package main

import (
	"context"
	"fmt"
	"io"
	"os"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

const QCOW_NAME = "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2"

var QCOW_LOCAL = fmt.Sprintf("/root/.panda/%s", QCOW_NAME)

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
	err = copyFileToContainerHelper(ctx, QCOW_LOCAL, QCOW_NAME, agent)
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
	ndl_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, recording.NdlogFilename())
	ndl_local, err := os.Create(ndl_dest)
	if err != nil {
		panic(err)
	}
	nBytes, err := io.Copy(ndl_local, ndl)
	if err != nil {
		panic(err)
	}
	if nBytes == 0 {
		panic("Bad copy")
	}
	defer ndl.Close()

	snp, err := recording.OpenSnapshot(ctx)
	if err != nil {
		panic(err)
	}
	snp_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, recording.SnapshotFilename())
	snp_local, err := os.Create(snp_dest)
	if err != nil {
		panic(err)
	}
	nBytes, err = io.Copy(snp_local, snp)
	if err != nil {
		panic(err)
	}
	if nBytes == 0 {
		panic("Bad copy")
	}
	defer snp.Close()

	fmt.Printf("Snapshot file: %s\n", recording.SnapshotFilename())
	fmt.Printf("Nondet log file: %s\n", recording.NdlogFilename())

	// Replay agent
	replay_agent, err := controller.CreateDockerPandaAgent2(ctx)
	if err != nil {
		panic(err)
	}
	defer replay_agent.Close()

	err = replay_agent.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = copyFileToContainerHelper(ctx, QCOW_LOCAL, QCOW_NAME, replay_agent)
	if err != nil {
		panic(err)
	}
	err = copyFileToContainerHelper(ctx, snp_dest, recording.SnapshotFilename(), replay_agent)
	if err != nil {
		panic(err)
	}
	err = copyFileToContainerHelper(ctx, ndl_dest, recording.NdlogFilename(), replay_agent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting replay")
	replay, err := replay_agent.StartReplay(ctx, "test")
	if err != nil {
		panic(err)
	}
	println(replay.Serial)
	println(replay.Replay)

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

	err = replay_agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}

// ctx - context
// srcFilePath - file path on local machine
// dstFileName - name of the file in the container
// agent - PandaAgent to container to copy into
func copyFileToContainerHelper(ctx context.Context, srcFilePath string, dstFilename string, agent *controller.DockerPandaAgent) error {
	fileReader, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return err
	}
	err = agent.CopyFileToContainer(ctx, fileReader, fileInfo.Size(), dstFilename)
	return err
}
