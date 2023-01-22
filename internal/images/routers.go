package images

import (
	"context"
	"net/http"

	"github.com/panda-re/panda_studio/internal/db"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type imagesController struct {
	imageRepo ImageRepository
}

func ImagesRegister(router *gin.RouterGroup) {
	imageRepo, err := GetRepository(context.TODO())
	if err != nil {
		panic(err)
	}
	controller := imagesController{
		imageRepo: imageRepo,
	}

	// List all Images
	router.GET("/", controller.ImageList)
	// List a single Image
	router.GET("/:imageId", controller.ImageGet)
	// Create a new file for an image
	router.POST("/:imageId/files", controller.ImageFileCreate)
	// Upload file to created file
	router.POST("/:imageId/files/:fileId/upload", controller.ImageFileUpload)

	// Create a new image
	router.POST("/", controller.ImageCreate)
	// delete an image
	router.DELETE("/:imageId", controller.ImageDelete)
}

func (ic *imagesController) ImageList(ctx *gin.Context) {
	// todo: allow search criteria
	images, err := ic.imageRepo.FindAll(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, images)
}

func (ic *imagesController) ImageGet(ctx *gin.Context) {
	id := ctx.Param("imageId")

	image, err := ic.imageRepo.FindOne(ctx, db.ParseObjectID(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, image)
}

func (*imagesController) ImageFile(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (ic *imagesController) ImageFileCreate(ctx *gin.Context) {
	var req ImageFileCreateRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	req.ImageID = db.ParseObjectID(ctx.Param("imageId"))

	file, err := ic.imageRepo.CreateImageFile(ctx, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, file)
}

func (ic *imagesController) ImageFileUpload(ctx *gin.Context) {
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

	imgFile, err := ic.imageRepo.UploadImageFile(ctx, &ImageFileUploadRequest{
		ImageId: db.ParseObjectID(ctx.Param("imageId")),
		FileId: db.ParseObjectID(ctx.Param("fileId")),
	}, fileReader)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, imgFile)
}

func (*imagesController) ImageCreate(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageDelete(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}
