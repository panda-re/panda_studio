package api

import (
	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
	"net/http"
)

func (s *PandaStudioServer) CreateRecording(ctx *gin.Context) {
	//TODO: do some stuff to actually make a recording
}

func (s *PandaStudioServer) FindRecordingById(ctx *gin.Context, recordingId string) {
	recording, err := s.recordingRepo.FindRecording(ctx, db.ParseObjectID(recordingId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.JSON(http.StatusOK, recording)
}

func (s *PandaStudioServer) DeleteRecordingById(ctx *gin.Context, recordingId string) {
	deleted, err := s.recordingRepo.DeleteRecording(ctx, db.ParseObjectID(recordingId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}
