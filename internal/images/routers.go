package images

import (
	"context"
	"net/http"

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
	router.GET("/:id", controller.ImageGet)
	// In the future, images may have more than one file attached
	router.GET("/:id/download", controller.ImageDownload)
	router.POST("/", controller.ImageCreate)
	router.DELETE("/:id", controller.ImageDelete)
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

func (*imagesController) ImageGet(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageDownload(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageCreate(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}

func (*imagesController) ImageDelete(ctx *gin.Context) {
	ctx.Error(errors.New("not implemented"))
}
