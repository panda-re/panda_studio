package main

import (
	"context"
	"time"

	"github.com/panda-re/panda_studio/executor"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	err := executor.RunCommand(ctx, "uname -a")
	if err != nil {
		panic(err)
	}
}
