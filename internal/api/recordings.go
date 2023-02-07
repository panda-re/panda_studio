package api

import (
	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
	"net/http"
)

func (s *PandaStudioServer) CreateRecording(ctx *gin.Context) {
	//TODO: add domain logic for creating recording
}

func (s *PandaStudioServer) FindAllRecordings(ctx *gin.Context) {
	recordings, err := s.recordingRepo.FindAllRecordings(ctx)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.JSON(http.StatusOK, recordings)
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

func (s *PandaStudioServer) CreateRecordingFile(ctx *gin.Context, recordingId string) {
	//TODO: add domain logic for creating recording
}

func (s *PandaStudioServer) DeleteRecordingFile(ctx *gin.Context, recordingId string, fileId string) {
	recordingFile, err := s.recordingRepo.DeleteRecordingFile(ctx, db.ParseObjectID(recordingId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, recordingFile)
}
