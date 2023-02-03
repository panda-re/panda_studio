package api

import (
	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
	"net/http"
)

func (s *PandaStudioServer) FindRecordingById(ctx *gin.Context, recordingId string) {
	image, err := s.recordingRepo.FindRecording(ctx, db.ParseObjectID(recordingId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.JSON(http.StatusOK, image)
}

func (s *PandaStudioServer) DeleteRecordingById(ctx *gin.Context, recordingId string) {
	deleted, err := s.recordingRepo.DeleteRecording(ctx, db.ParseObjectID(recordingId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}
