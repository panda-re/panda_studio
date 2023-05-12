package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

const QCOW_NAME = "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2"
const RECORDING_NAME = "test"
const EXECUTION_LOG = "execution.log"

var QCOW_LOCAL = fmt.Sprintf("/root/.panda/%s", QCOW_NAME)

func main() {
	ctx := context.Background()

	runRecording(ctx)

	runReplay(ctx)
}

func runRecording(ctx context.Context) {
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

	ndl_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, recording.NdlogFilename())
	err = copyFileFromContainerHelper(ctx, recording.NdlogFilename(), ndl_dest, agent)
	if err != nil {
		panic(err)
	}
	snp_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, recording.SnapshotFilename())
	err = copyFileFromContainerHelper(ctx, recording.SnapshotFilename(), snp_dest, agent)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Snapshot file: %s\n", recording.SnapshotFilename())
	fmt.Printf("Nondet log file: %s\n", recording.NdlogFilename())

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}

func runReplay(ctx context.Context) {
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
	snp_name := fmt.Sprintf("%s-rr-snp", RECORDING_NAME)
	snp_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, snp_name)
	err = copyFileToContainerHelper(ctx, snp_dest, snp_name, replay_agent)
	if err != nil {
		panic(err)
	}
	ndl_name := fmt.Sprintf("%s-rr-nondet.log", RECORDING_NAME)
	ndl_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, ndl_name)
	err = copyFileToContainerHelper(ctx, ndl_dest, ndl_name, replay_agent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting replay")
	replay, err := replay_agent.StartReplay(ctx, RECORDING_NAME)
	if err != nil {
		panic(err)
	}
	println(replay.Serial)
	println(replay.Replay)
	// Uncomment to get execution log from agent
	// For more debugging log output, see /docker/Dockerfile.panda-agent
	// log_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, EXECUTION_LOG)
	// copyFileFromContainerHelper(ctx, EXECUTION_LOG, log_dest, replay_agent)

	err = replay_agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}

// ctx - context
// srcFilePath - file path on local machine
// dstFilePath - name of the file in the container
// agent - PandaAgent to container to copy into
func copyFileToContainerHelper(ctx context.Context, srcFilePath string, dstFilePath string, agent *controller.DockerPandaAgent) error {
	fileReader, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return err
	}
	err = agent.CopyFileToContainer(ctx, fileReader, fileInfo.Size(), dstFilePath)
	return err
}

// ctx - context
// srcFilePath - file path in container to copy from
// dstFilePath - file path on local machine to copy to
// agent - PandaAgent to container to copy from
func copyFileFromContainerHelper(ctx context.Context, srcFilePath string, dstFilePath string, agent *controller.DockerPandaAgent) error {
	src, err := agent.CopyFileFromContainer(ctx, srcFilePath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	nBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	if nBytes == 0 {
		return errors.New("did not copy file")
	}
	return nil
}
