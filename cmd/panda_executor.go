package main

import (
	"context"
	"time"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	err := controller.RunCommand(ctx, "uname -a")
	if err != nil {
		panic(err)
	}
}
