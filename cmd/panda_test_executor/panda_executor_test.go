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
	t.Cleanup(func() {
		fmt.Printf("Number of tests: %d\nNumber passed: %d\nSuccess rate: %d%%\n", num_tests, num_passed, 100*num_passed/num_tests)
	})
	// TODO tests
	t.Run("Agent", TestAgent)
	if t.Failed() {
		t.Fatal("Agent unsuccessful")
	}
	t.Run("Recording", TestRecord)
	if t.Failed() {
		t.Fatal("Recording unsuccessful")
	}
	t.Run("Replay", TestReplay)
}

var agent controller.PandaAgent

var recording_name string = "test"
var ctx = context.Background()

func TestAgent(t *testing.T) {
	var err error
	t.Cleanup(func() {
		err = agent.StopAgent(ctx)
		if err != nil {
			t.Fatal(err)
		}
		err = agent.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
	agent, err = controller.CreateDefaultDockerPandaAgent(ctx, "/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("PreCommand", TestPrematureCommand)
	if !t.Failed() {
		num_passed++
	}
	t.Run("PreStop", TestPrematureStop)
	if !t.Failed() {
		num_passed++
	}
	err = agent.StartAgent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("ExtraStart", TestExtraStart)
	if !t.Failed() {
		num_passed++
	}
	t.Run("Commands", TestCommands)
	if !t.Failed() {
		num_passed++
	}
}
func TestPrematureCommand(t *testing.T) {
	num_tests++
	_, err := agent.RunCommand(ctx, "")
	if err == nil {
		t.Error("Did not prevent premature command")
	}
}

func TestPrematureStop(t *testing.T) {
	num_tests++
	err := agent.StopAgent(ctx)
	if err == nil {
		t.Error("Did not prevent premature stop")
	}
}

func TestExtraStart(t *testing.T) {
	num_tests++
	if err := agent.StartAgent(ctx); err == nil {
		t.Error("Did not prevent a second PANDA start")
	}
}

func TestCommands(t *testing.T) {
	num_tests++
	commands := []string{
		"echo Hello World",
	}
	for _, cmd := range commands {
		response, err := agent.RunCommand(ctx, cmd)
		if err != nil {
			t.Error(err)
		}
		if response == nil {
			t.Fatal("Did not receive a response")
		} else if response.Logs != "Hello World" {
			t.Fatal("Did not receive correct response from command")
		}
	}
}

func TestRecord(t *testing.T) {
	var err error
	t.Cleanup(func() {
		err = agent.StopAgent(ctx)
		if err != nil {
			t.Fatal(err)
		}
		err = agent.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
	agent, err = controller.CreateDefaultDockerPandaAgent(ctx, "/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2")
	if err != nil {
		t.Fatal(err)
	}

	err = agent.StartAgent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// TODO more and proper tests
	t.Run("Start", TestStartRecording)
	if t.Failed() {
		t.FailNow()
	} else {
		num_passed++
	}
	t.Run("Commands", TestCommands)
	if !t.Failed() {
		num_passed++
	}
}

func TestStartRecording(t *testing.T) {
	num_tests++
	if err := agent.StartRecording(ctx, recording_name); err != nil {
		t.Error(err)
	}
}

var replay_agent controller.PandaReplayAgent

func TestReplay(t *testing.T) {
	var err error
	t.Cleanup(func() {
		err = replay_agent.StopAgent(ctx)
		if err != nil {
			t.Fatal(err)
		}
		err = replay_agent.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	replay_agent, err = controller.CreateReplayDockerPandaAgent(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// TODO more and proper tests

	t.Run("PreStop", TestPrematureReplayStop)
	if !t.Failed() {
		num_passed++
	}
	t.Run("RunReplay", TestRunReplay)
	if !t.Failed() {
		num_passed++
	}
}

func TestPrematureReplayStop(t *testing.T) {
	num_tests++
	_, err := replay_agent.StopReplay(ctx)
	if err == nil {
		t.Error("Did not prevent premature stop")
	}
}

func TestRunReplay(t *testing.T) {
	num_tests++
	replay, err := replay_agent.StartReplayAgent(ctx, recording_name)
	if err != nil {
		t.Fatal(err)
	}
	if replay.Serial == "" || replay.Replay == "" {
		// TODO better job ensuring the logs are correct
		t.Error("Test Failed. Replay returned partially incomplete")
	}
}
