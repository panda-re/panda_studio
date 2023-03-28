package api

import (
	"context"

	"github.com/panda-re/panda_studio/internal/db/repos"
	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/panda_studio.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/panda_studio.yaml

type PandaStudioServer struct {
	imageRepo     repos.ImageRepository
	recordingRepo repos.RecordingRepository
	programRepo  repos.ProgramRepository
	programExecutor *controller.PandaProgramExecutor
}

// Ensures our implementation conforms to the correct interface
var _ ServerInterface = &PandaStudioServer{}

func NewPandaStudioServer() (*PandaStudioServer, error) {
	ctx := context.TODO()
	imageRepo, err := repos.GetImageRepository(ctx)
	if err != nil {
		return nil, err
	}

	recordingRepo, err := repos.GetRecordingRepository(ctx)
	if err != nil {
		return nil, err
	}

	programRepo, err := repos.GetProgramRepository(ctx)
	if err != nil {
		return nil, err
	}

	programExecutor, err := controller.NewPandaProgramExecutor(ctx)
	if err != nil {
		return nil, err
	}

	return &PandaStudioServer{
		imageRepo:     imageRepo,
		recordingRepo: recordingRepo,
		programRepo:  programRepo,
		programExecutor: programExecutor,
	}, nil
}
