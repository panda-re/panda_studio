package api

import (
	"context"

	"github.com/panda-re/panda_studio/internal/images"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/panda_studio.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/panda_studio.yaml

type PandaStudioServer struct {
	imageRepo images.ImageRepository
}

// Ensures our implementation conforms to the correct interface
var _ ServerInterface = &PandaStudioServer{}

func NewPandaStudioServer() (*PandaStudioServer, error) {
	imageRepo, err := images.GetRepository(context.TODO())
	if err != nil {
		return nil, err
	}

	return &PandaStudioServer{
		imageRepo: imageRepo,
	}, nil
}
