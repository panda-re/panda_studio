package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/panda_controller"
	"github.com/panda-re/panda_studio/panda_agent/pb"
	"github.com/pkg/errors"
)

func (s *PandaStudioServer) FindAllImages(ctx *gin.Context) {
	// todo: allow search criteria
	images, err := s.imageRepo.FindAll(ctx)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, images)
}

func (s *PandaStudioServer) FindImageById(ctx *gin.Context, imageId string) {
	image, err := s.imageRepo.FindOne(ctx, db.ParseObjectID(imageId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, image)
}

func (s *PandaStudioServer) CreateImage(ctx *gin.Context) {
	var createReq CreateImageRequest
	err := ctx.BindJSON(&createReq)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	var newConfig pb.PandaConfig
	temporaryVariable, _ := json.Marshal(createReq.Config)
	err = json.Unmarshal(temporaryVariable, &newConfig)
	if err != nil {
		ctx.Error(errors.Wrap(err, "unable to recast config"))
		return
	}

	created, err := s.imageRepo.Create(ctx, &models.Image{
		Name:        *createReq.Name,
		Description: *createReq.Description,
		Config: &models.ImageConfiguration{
			PandaConfig: newConfig,
		},
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (s *PandaStudioServer) DeleteImageById(ctx *gin.Context, imageId string) {
	deleted, err := s.imageRepo.DeleteOne(ctx, db.ParseObjectID(imageId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}

func (s *PandaStudioServer) UpdateImage(ctx *gin.Context, imageId string) {
	var updateReq CreateImageRequest
	err := ctx.BindJSON(&updateReq)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	var newConfig pb.PandaConfig
	temporaryVariable, _ := json.Marshal(updateReq.Config)
	err = json.Unmarshal(temporaryVariable, &newConfig)
	if err != nil {
		ctx.Error(errors.Wrap(err, "unable to recast config"))
		return
	}

	updated, err := s.imageRepo.Update(ctx, &models.Image{
		Name:        *updateReq.Name,
		Description: *updateReq.Description,
		ID:          db.ParseObjectID(imageId),
		Config: &models.ImageConfiguration{
			PandaConfig: newConfig,
		},
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, updated)
}

func (s *PandaStudioServer) CreateImageFile(ctx *gin.Context, imageId string) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileObj, err := s.imageRepo.CreateImageFile(ctx, &models.ImageFileCreateRequest{
		ImageID:  db.ParseObjectID(imageId),
		FileName: form.Value["file_name"][0],
		FileType: form.Value["file_type"][0],
	})
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileReader, err := formFile.Open()
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer fileReader.Close()

	fileObj, err = s.imageRepo.UploadImageFile(ctx, &models.ImageFileUploadRequest{
		ImageId: db.ParseObjectID(imageId),
		FileId:  fileObj.ID,
	}, fileReader)

	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, fileObj)
}

func (s *PandaStudioServer) CreateImageFileFromUrl(ctx *gin.Context, imageId string) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileObj, err := s.imageRepo.CreateImageFile(ctx, &models.ImageFileCreateRequest{
		ImageID:  db.ParseObjectID(imageId),
		FileName: form.Value["file_name"][0],
		FileType: form.Value["file_type"][0],
	})
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	resp, err := http.Get(form.Value["url"][0])
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode)
		ctx.Error(errors.WithStack(err))
		return
	}

	fileObj, err = s.imageRepo.UploadImageFile(ctx, &models.ImageFileUploadRequest{
		ImageId: db.ParseObjectID(imageId),
		FileId:  fileObj.ID,
	}, resp.Body)

	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, fileObj)
}

func (s *PandaStudioServer) DownloadImageFile(ctx *gin.Context, imageId ImageId, fileId FileId) {
	imgFile, err := s.imageRepo.FindOneImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}

	fileReader, err := s.imageRepo.OpenImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}
	defer fileReader.Close()

	ctx.DataFromReader(http.StatusOK, imgFile.Size, "application/octet-stream", fileReader, map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename=%s`, imgFile.FileName),
	})
}

func (s *PandaStudioServer) DeleteImageFile(ctx *gin.Context, imageId ImageId, fileId FileId) {
	imgFile, err := s.imageRepo.DeleteImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, imgFile)
}

func (s *PandaStudioServer) CreateDerivedImage(ctx *gin.Context, imageId string) {
	var deriveReq DeriveImageFileRequest
	err := ctx.BindJSON(&deriveReq)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	oldImageFile, err := s.imageRepo.FindOne(ctx, db.ParseObjectID(imageId))
	if err != nil {
		ctx.Error(err)
		return
	}

	diExecutor, err := s.deriveImageJob.NewDeriveImageJobExecutor(ctx, &panda_controller.DeriveImageJobParams{
		BaseImage:   oldImageFile,
		NewImage:    *deriveReq.Newname,
		Resize:      *deriveReq.Size, //TODO change to string in DeriveImageFileRequest
		DockerImage: *deriveReq.Dockerhubimagename,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	err = diExecutor.Run(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
}
