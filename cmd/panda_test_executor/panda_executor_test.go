package main

import (
	"context"
	"fmt"
	"testing"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

var num_passed int = 0
var num_tests int = 0

func TestMain(t *testing.T) {
	// TODO tests
	t.Run("Recording", TestRecord)
	t.Run("Replay", TestReplay)

	fmt.Printf("Number of tests: %d\nNumber passed: %d\nSuccess rate: %d%%\n", num_tests, num_passed, 100*num_passed/num_tests)
}

var agent controller.PandaAgent

// var recording_name string = "test"
var ctx = context.Background()

func TestRecord(t *testing.T) {
	var err error
	agent, err = controller.CreateDefaultDockerPandaAgent(ctx, "/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		panic(err)
	}
	defer agent.Close()

	// TODO more and proper tests
	t.Run("PreCommand", TestPrematureCommand)

	err = agent.StartAgent(ctx)
	if err != nil {
		panic(err)
	}

	err = agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
}

func TestPrematureCommand(t *testing.T) {
	num_tests++
	_, err := agent.RunCommand(ctx, "")
	if err != nil {
		num_passed++
	} else {
		t.Error("Did not prevent premature command")
	}
}

var replay_agent controller.PandaReplayAgent

func TestReplay(t *testing.T) {
	var err error
	replay_agent, err = controller.CreateReplayDockerPandaAgent(ctx)
	if err != nil {
		panic(err)
	}
	defer replay_agent.Close()
	if err != nil {
		panic(err)
	}
	// TODO more and proper tests

	t.Run("PreStop", TestPrematureStop)

	// err = replay_agent.StopAgent(ctx)
	// if err != nil {
	// 	panic(err)
	// }
}

func TestPrematureStop(t *testing.T) {
	num_tests++
	_, err := replay_agent.StopReplay(ctx)
	if err != nil {
		num_passed++
	} else {
		t.Error("Did not prevent premature stop")
	}
}
