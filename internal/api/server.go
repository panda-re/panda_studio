package api

import (
	"context"

	"github.com/panda-re/panda_studio/internal/db/repos"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/panda_studio.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/panda_studio.yaml

type PandaStudioServer struct {
	imageRepo repos.ImageRepository
}

// Ensures our implementation conforms to the correct interface
var _ ServerInterface = &PandaStudioServer{}

func NewPandaStudioServer() (*PandaStudioServer, error) {
	imageRepo, err := repos.GetImageRepository(context.TODO())
	if err != nil {
		return nil, err
	}

	return &PandaStudioServer{
		imageRepo: imageRepo,
	}, nil
}
