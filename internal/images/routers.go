package images

import (
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type imagesController struct {}

func ImagesRegister(router *gin.RouterGroup) {
	controller := imagesController{}

	router.GET("/", controller.ImageList)
	router.GET("/:id", controller.ImageGet)
	// In the future, images may have more than one file attached
	router.GET("/:id/download", controller.ImageDownload)
	router.POST("/", controller.ImageCreate)
	router.DELETE("/:id", controller.ImageDelete)
}

func (*imagesController) ImageList(c *gin.Context) {
	// c.JSON(http.StatusInternalServerError, errors.New("not implemented"))
	c.Error(errors.New("not implemented"))
}

func (*imagesController) ImageGet(c *gin.Context) {
	c.Error(errors.New("not implemented"))
}

func (*imagesController) ImageDownload(c *gin.Context) {
	c.Error(errors.New("not implemented"))
}

func (*imagesController) ImageCreate(c *gin.Context) {
	c.Error(errors.New("not implemented"))
}

func (*imagesController) ImageDelete(c *gin.Context) {
	c.Error(errors.New("not implemented"))
}
