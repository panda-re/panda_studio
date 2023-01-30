package main

import (
	"context"
	"fmt"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

func main() {
	ctx := context.Background()

	agent, err := controller.CreateReplayDockerPandaAgent(ctx)
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	if err != nil {
		panic(err)
	}
	// Do not use StartAgent. Starts PANDA, which will prevent replay
	// fmt.Println("Starting agent")
	// err = agent.StartAgent(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Starting replay")
	replay, err := agent.StartReplay(ctx, "test")
	if err != nil {
		panic(err)
	}
	println(replay.Serial)
	println(replay.Replay)

	// Replay automatically stops after being run, only for cancelling
	// replay, err = agent.StopReplay(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}
