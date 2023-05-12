package panda_controller

import (
	"context"
	"fmt"
	"io"

	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/db/repos"
	"github.com/pkg/errors"
)

type DeriveImageJob struct {
	imageRepo repos.ImageRepository
}

type DeriveImageJobExecutor struct {
	imageRepo repos.ImageRepository
	params    DeriveImageJobParams
	agent     *dockerPandaDIAgent
}

type DeriveImageJobParams struct {
	BaseImage   *models.Image
	NewImage    string
	Resize      string
	DockerImage string
}

func NewDeriveImageJob(ctx context.Context) (*DeriveImageJob, error) {
	imageRepo, err := repos.GetImageRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &DeriveImageJob{
		imageRepo: imageRepo,
	}, nil
}

func (di *DeriveImageJob) NewDeriveImageJobExecutor(ctx context.Context, params *DeriveImageJobParams) (*DeriveImageJobExecutor, error) {
	diJobExecutor := &DeriveImageJobExecutor{
		imageRepo: di.imageRepo,
		params:    *params,
	}
	return diJobExecutor, nil
}

// setupContainer
func (di *DeriveImageJobExecutor) setupContainer(ctx context.Context) error {
	//download the image
	image := di.params.BaseImage
	file := &models.ImageFile{}

	for _, curFile := range image.Files {
		if curFile.FileType == "qcow2" {
			file = curFile
		}
	}
	if file.FileName == "" {
		return errors.Wrap(nil, "qcow image could not be found")
	}

	fileReader, err := di.imageRepo.OpenImageFile(ctx, image.ID, file.ID)
	if err != nil {
		return errors.Wrap(err, "failed to open image")
	}

	diAgent, err := CreateDockerPandaDIAgent(ctx, PandaDIAgentParams{
		base_image_name: di.params.BaseImage.Name,
		new_image:       di.params.NewImage,
		resize:          di.params.Resize,
		docker_image:    di.params.DockerImage,
		base_image:      fileReader,
		base_image_size: file.Size,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create the derive image container")
	}
	di.agent = diAgent

	return nil
}

// Run
func (di *DeriveImageJobExecutor) Run(ctx context.Context) error {
	//setup the derive imag container
	err := di.setupContainer(ctx)
	if err != nil {
		return err
	}
	defer di.agent.Close()

	//run the derive image job in the derive image agent
	err = di.agent.Run(ctx)
	if err != nil {
		fmt.Printf("Error running the derive image job: %s\n", err.Error())
		return err
	}

	//get derived image from container
	imageFile, err := di.agent.CopyFileFromContainer(ctx, di.params.NewImage) //TODO: verify
	if err != nil {
		return err
	}

	//upload image to object storage
	err = di.uploadImage(ctx, imageFile)
	if err != nil {
		return err
	}

	return nil
}

// uploadImage
func (di *DeriveImageJobExecutor) uploadImage(ctx context.Context, createdImageFile io.Reader) error {
	created, err := di.imageRepo.Create(ctx, &models.Image{
		Name:        di.params.NewImage,
		Description: fmt.Sprintf("Derived from %s", di.params.BaseImage.Name),
	})
	if err != nil {
		return err
	}

	imageFile, err := di.imageRepo.CreateImageFile(ctx, &models.ImageFileCreateRequest{
		ImageID:  created.ID,
		FileName: created.Name,
		FileType: "qcow2",
	})
	if err != nil {
		return err
	}

	fileObj, err := di.imageRepo.UploadImageFile(ctx, &models.ImageFileUploadRequest{
		ImageId: created.ID,
		FileId:  imageFile.ID,
	}, createdImageFile)
	if err != nil || fileObj == nil {
		return err
	}

	return nil
}
