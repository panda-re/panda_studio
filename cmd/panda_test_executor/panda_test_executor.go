package main

import (
	"context"
	"fmt"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

func main() {
	// Variables for keeping track of pass/fail of test
	num_tests := 0
	num_passed := 0
	recording_name := "test"

	ctx := context.Background()
	agent, err := controller.CreateDefaultDockerPandaAgent(ctx, "")
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	if err != nil {
		panic(err)
	}

	commands := []string{
		"uname -a",
		"ls /",
		"touch /NEW_FILE.txt",
		"ls /",
	}

	println("Testing PANDA Agent")
	// Running a command before PANDA starts
	num_tests++
	for _, cmd := range commands {
		_, err := agent.RunCommand(ctx, cmd)
		if err != nil {
			println("Test Passed. Prevented running command before PANDA starts")
			num_passed++
			break
		} else {
			println("Test Failed. Command ran without starting PANDA")
			panic(err)
		}
	}

	// Stopping PANDA before starting
	num_tests++
	err = agent.StopAgent(ctx)
	if err != nil {
		println("Test Passed. Prevented stopping PANDA before PANDA starts")
		num_passed++
	} else {
		println("Test Failed. Stopped nonexistent PANDA")
		panic(err)
	}

	fmt.Println("Starting agent")
	err = agent.StartAgent(ctx)
	if err != nil {
		panic(err)
	}

	// Starting second PANDA agent
	num_tests++
	err = agent.StartAgent(ctx)
	if err != nil {
		fmt.Println("Test Passed. Prevented second PANDA agent creation")
		num_passed++
	} else {
		fmt.Println("Test Failed. Created second PANDA agent creation")
		panic(err)
	}

	// Stopping recording before starting
	num_tests++
	_, err = agent.StopRecording(ctx)
	if err != nil {
		println("Test Passed. Prevented stopping of a non-existent recording")
		num_passed++
	} else {
		println("Test Failed. Stopped nothing from ever recording")
		panic(err)
	}

	if err := agent.StartRecording(ctx, recording_name); err != nil {
		panic(err)
	}

	// Starting second recording while one is in progress
	num_tests++
	if err := agent.StartRecording(ctx, "testing"); err != nil {
		println("Test Passed. Prevented starting a second concurrent recording")
		num_passed++
	} else {
		println("Test Failed. Started a second recording")
		panic(err)
	}

	for _, cmd := range commands {
		_, err = agent.RunCommand(ctx, cmd)
		if err != nil {
			panic(err)
		}
	}

	_, err = agent.StopRecording(ctx)
	if err != nil {
		panic(err)
	}

	replay_agent, err := controller.CreateReplayDockerPandaAgent(ctx)
	if err != nil {
		panic(err)
	}
	defer replay_agent.Close()
	if err != nil {
		panic(err)
	}
	println("Testing PANDA replay")
	// Testing stopping a replay before starting
	num_tests++
	if _, err := replay_agent.StopReplay(ctx); err != nil {
		println("Test Passed. Prevented stopping of a non-existent replay")
		num_passed++
	} else {
		println("Test Failed. Stopped nothing from being replayed")
	}

	// Starting a replay of a recording that doesn't exist
	num_tests++
	if _, err := replay_agent.StartReplayAgent(ctx, "DNE"); err != nil {
		println("Test Passed. Prevented replay of a recording that doesn't exist")
		num_passed++
	} else {
		println("Test Failed. Attempted to replay a nonexistant recording")
	}

	// Testing return of serial and execution
	num_tests++
	replay, err := replay_agent.StartReplayAgent(ctx, recording_name)
	if err != nil {
		println("Test Failed. Could not replay")
		panic(err)
	}
	if replay.Serial != "" && replay.Replay != "" {
		// TODO better job ensuring the logs are correct
		println("Test Passed. Replay ran successfully")
		num_passed++
	} else {
		println("Test Failed. Replay returned partially incomplete")
	}

	fmt.Printf("Number of tests: %d\nNumber passed: %d\nSuccess rate: %d%%\n", num_tests, num_passed, 100*num_passed/num_tests)
	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
	err = replay_agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}
