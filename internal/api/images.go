package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/images"
)

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

func (s *PandaStudioServer) FindImageById(ctx *gin.Context, imageId string) {
	image, err := s.imageRepo.FindOne(ctx, db.ParseObjectID(imageId))
	if err != nil {
		ctx.Error(err)
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, image)
}

func (s *PandaStudioServer) CreateImageFile(ctx *gin.Context, imageId string) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(err)
		return
	}
	
	fileObj, err := s.imageRepo.CreateImageFile(ctx, &images.ImageFileCreateRequest{
		ImageID: db.ParseObjectID(imageId),
		FileName: form.Value["file_name"][0],
		FileType: form.Value["file_type"][0],
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}

	fileReader, err := formFile.Open()
	if err != nil {
		ctx.Error(err)
		return
	}
	defer fileReader.Close()

	fileObj, err = s.imageRepo.UploadImageFile(ctx, &images.ImageFileUploadRequest{
		ImageId: db.ParseObjectID(imageId),
		FileId: fileObj.ID,
	}, fileReader)

	if err != nil {
		ctx.Error(err)
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, fileObj)
}