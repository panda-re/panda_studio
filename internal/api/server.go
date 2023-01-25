package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
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

func (s *PandaStudioServer) FindAllImages(ctx *gin.Context) {
	// todo: allow search criteria
	images, err := s.imageRepo.FindAll(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, images)
}

func (s *PandaStudioServer) FindImageById(ctx *gin.Context, id string) {
	image, err := s.imageRepo.FindOne(ctx, db.ParseObjectID(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, image)
}