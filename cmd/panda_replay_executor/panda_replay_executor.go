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
	// Necessary because regular printing won't do newlines
	// Offset because of log overheads
	length := len(replay.Logs) - 1
	for i := 8; i < length; i++ {
		// Check for \r and \n
		c := replay.Logs[i]
		if c == '\\' && i+1 < length {
			switch char := replay.Logs[i+1]; char {
			case 'r':
				fmt.Printf("\r")
				i++
				continue
			case 'n':
				fmt.Printf("\n")
				i++
				continue
			case 't':
				fmt.Printf("\t")
				i++
				continue
			}
		}
		fmt.Printf("%c", replay.Logs[i])
	}
	println()
	// fmt.Printf("%s\n", replay.Logs)

	// Replay automatically stops after being run, only for cancelling
	// err = agent.StopReplay(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// time.Sleep(time.Second * 20)

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}
