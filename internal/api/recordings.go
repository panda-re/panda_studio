package api

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"
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

func (s *PandaStudioServer) DownloadRecordingFile(ctx *gin.Context, recordingId ImageId, fileId FileId) {
	recordingFile, err := s.recordingRepo.FindRecordingFile(ctx, db.ParseObjectID(recordingId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}

	fileReader, err := s.recordingRepo.OpenRecordingFile(ctx, db.ParseObjectID(recordingId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}
	defer fileReader.Close()

	ctx.DataFromReader(http.StatusOK, recordingFile.Size, "application/octet-stream", fileReader, map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename=%s`, recordingFile.Name),
	})
}
