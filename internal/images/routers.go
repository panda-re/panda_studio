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

	router.GET("/", controller.ImageList)
	router.GET("/:imageId", controller.ImageGet)
	// In the future, images may have more than one file attached
	router.GET("/:imageId/files/:fileId", controller.ImageFile)
	router.POST("/", controller.ImageCreate)
	router.DELETE("/imageId", controller.ImageDelete)
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

	image, err := ic.imageRepo.FindOne(ctx, db.ObjectID(id))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, image)
}

func (*imagesController) ImageFile(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageCreate(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageUploadFile(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageDelete(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}
