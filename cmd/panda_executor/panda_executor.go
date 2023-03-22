package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	config "github.com/panda-re/panda_studio/internal/configuration"
	"github.com/panda-re/panda_studio/internal/db/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

//go:embed test_program.txt
var testProgram string

//go:embed bionic-image.json
var testImageSpec string

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
	IMAGE_FILE := "images/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2"
	// IMAGE_FILE := "images/ubuntu-2204.qcow2"
	file, err := os.Open(IMAGE_FILE)
	if err != nil {
		panic(err)
	}

	fileSize, err := getReaderSize(file)
	if err != nil {
		panic(err)
	}

	var image models.Image
	if err := json.Unmarshal([]byte(testImageSpec), &image); err != nil {
		panic(err)
	}

	// debug print the image
	fmt.Printf("image: %+v\n", image)

	jobOpts := controller.PandaProgramExecutorOptions{
		Image: &image,
		Program: &models.InteractionProgram{
			ID: &primitive.NilObjectID,
			Name: "test_program",
			Instructions: testProgram,
		},
		Instructions: prog,
		ImageFileReader: file,
		ImageFileSize: fileSize,
	}

	job, err := progExec.NewExecutorJob(ctx, &jobOpts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("job: %v\n", job)

	job.StartJob(ctx)
}

func getReaderSize(reader io.ReadSeeker) (size int64, err error) {
	// Get the size of the file
	// https://stackoverflow.com/a/24563853
	currentPos, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	defer reader.Seek(currentPos, io.SeekStart)

	endPos, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	defer reader.Seek(currentPos, io.SeekStart)

	return endPos, nil
}
