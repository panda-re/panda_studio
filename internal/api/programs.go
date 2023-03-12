package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/pkg/errors"
)

// FindAllPrograms implements ServerInterface
func (s *PandaStudioServer) FindAllPrograms(ctx *gin.Context) {
	// todo: allow search criteria
	programs, err := s.programRepo.FindAll(ctx)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.JSON(http.StatusOK, programs)
}

// FindProgramById implements ServerInterface
func (s *PandaStudioServer) FindProgramById(ctx *gin.Context, programId string) {
	program, err := s.programRepo.FindOne(ctx, db.ParseObjectID(programId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.JSON(http.StatusOK, program)
}

// CreateProgram implements ServerInterface
func (s *PandaStudioServer) CreateProgram(ctx *gin.Context) {
	var toCreate models.InteractionProgram
	err := ctx.BindJSON(&toCreate)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	created, err := s.programRepo.Create(ctx, &toCreate)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

// UpdateProgramById implements ServerInterface
func (s *PandaStudioServer) UpdateProgramById(ctx *gin.Context, programId string) {
	var toUpdate models.InteractionProgram
	err := ctx.BindJSON(&toUpdate)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	updated, err := s.programRepo.Update(ctx, db.ParseObjectID(programId), &toUpdate)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// DeleteProgramById implements ServerInterface
func (s *PandaStudioServer) DeleteProgramById(ctx *gin.Context, programId string) {
	deleted, err := s.programRepo.Delete(ctx, db.ParseObjectID(programId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}
