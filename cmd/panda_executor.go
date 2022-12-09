package main

import (
	"github.com/panda-re/panda_studio/executor"
)

func main() {
	err := executor.RunDocker()
	if err != nil {
		panic(err)
	}
}
