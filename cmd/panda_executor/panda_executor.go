package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	config "github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

//go:embed test_program.txt
var testProgram string

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	fmt.Println(testProgram)

	prog, err := models.ParseInteractionProgram(testProgram)
	if err != nil {
		panic(err)
	}

	// debug print each item inprog
	fmt.Println("Instructions:")
	for _, item := range prog {
		// get the type of the item
		fmt.Printf("%T %+v\n", item, item)
	}

	progExec := controller.PandaProgramExecutor{}

	// open a stream to the file in blob storage
	file, err := os.Open("images/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		panic(err)
	}

	jobOpts := controller.PandaProgramExecutorOptions{
		Image: &models.Image{
			ID: &primitive.NilObjectID,
			Name: "default_image",
			Description: "Default Image",
			Files: []*models.ImageFile{
				{
					ID: &primitive.NilObjectID,
					ImageID: &primitive.NilObjectID,
					FileName: "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2",
					FileType: "qcow2",
					IsUploaded: true,
					Size: 0,
					Sha256: "",
				},
			},
			Config: &models.ImageConfiguration{},
		},
		Program: &models.InteractionProgram{
			ID: &primitive.NilObjectID,
			Name: "test_program",
			Instructions: testProgram,
		},
		Instructions: prog,
		ImageFileReader: file,
	}

	job, err := progExec.NewExecutorJob(&jobOpts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("job: %v\n", job)

	job.StartJob(ctx)
}

func old_main() {
	// Default agent
	ctx := context.Background()
	file, err := os.Open("../images/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		panic(err)
	}
	agent, err := controller.CreateDefaultDockerPandaAgent(ctx, file)
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
