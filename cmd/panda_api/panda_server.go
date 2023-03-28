package main

import (
	_ "embed"

	config "github.com/panda-re/panda_studio/internal/configuration"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/api"
	"github.com/panda-re/panda_studio/internal/middleware"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	if err := runServer(); err != nil {
		panic(err)
	}
}

func runServer() error {
	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	swagger, err := api.GetSwagger()
	if err != nil {
		return err
	}
	swagger.Servers = nil
	server, err := api.NewPandaStudioServer()
	if err != nil {
		return err
	}

	// r.Use(oapimiddleware.OapiRequestValidator(swagger))
	api.RegisterHandlersWithOptions(r, server, api.GinServerOptions{
		BaseURL: "/api",
	})

	if err := r.Run(); err != nil {
		return err
	}

	return nil
}