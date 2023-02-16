package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
)

// Each enum represents the state panda was in that caused the exception
// Enum should match that in /panda_agent/agent.py
const (
	RUNNING       = 0
	NOT_RUNNING   = 1
	RECORDING     = 2
	NOT_RECORDING = 3
	REPLAYING     = 4
	NOT_REPLAYING = 5
)

var error_to_string [6]string = [6]string{"RUNNING", "NOT_RUNNING", "RECORDING", "NOT_RECORDING", "REPLAYING", "NOT_REPLAYING"}

// Extracts the error number from an err from agent
func getError(err error) int {
	// Find the numbers in the error message
	re := regexp.MustCompile("[0-9]+")
	// Return the first one as an integer
	nums := re.FindAllString(err.Error(), -1)
	num, _ := strconv.Atoi(nums[0])
	return num
}

var num_passed int = 0
var num_tests int = 0

// Runs a test for the agent, recording and replay
// Prints the number of tests and success rate
func TestMain(t *testing.T) {
	t.Cleanup(func() {
		fmt.Printf("Number of tests: %d\nNumber passed: %d\nSuccess rate: %d%%\n", num_tests, num_passed, 100*num_passed/num_tests)
	})
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

// Recording name for testing record and replay
var recording_name string = "panda_executor_test"

// Consistent commands to run
var commands = []string{
	"echo Hello World",
	"uname -a",
	"touch /NEW_FILE.txt",
}

// Outputs of commands to ensure they ran correctly
var commands_output = []string{
	"Hello World",
	"Linux ubuntu 4.15.0-72-generic #81-Ubuntu SMP Tue Nov 26 12:20:02 UTC 2019 x86_64 x86_64 x86_64 GNU/Linux",
	"",
}
var ctx = context.Background()

// Tests to ensure the agent functions properly
// Tests premature execution and that commands return properly
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

// Tests executing a command before the agent has started
// The agent should prevent this from happening with an exception
// Should be run before agent.StartAgent
func TestPrematureCommand(t *testing.T) {
	num_tests++
	_, err := agent.RunCommand(ctx, "")
	if err == nil {
		t.Fatal("Did not prevent premature command")
	} else {
		err_num := getError(err)
		err_expected := NOT_RUNNING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Tests attempting to stop the agent before the agent has started
// The agent should prevent this from happening with an exception
// Should be run before agent.StartAgent
func TestPrematureStop(t *testing.T) {
	num_tests++
	err := agent.StopAgent(ctx)
	if err == nil {
		t.Fatal("Did not prevent premature stop")
	} else {
		err_num := getError(err)
		err_expected := NOT_RUNNING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Tests attempting to start the agent after the agent has started
// The agent should prevent this from happening with an exception
// Should be run after agent.StartAgent
func TestExtraStart(t *testing.T) {
	num_tests++
	err := agent.StartAgent(ctx)
	if err == nil {
		t.Fatal("Did not prevent a second PANDA start")
	} else {
		err_num := getError(err)
		err_expected := RUNNING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Tests sending serial commands
// Should be run after agent.StartAgent
func TestCommands(t *testing.T) {
	num_tests++
	for i, cmd := range commands {
		response, err := agent.RunCommand(ctx, cmd)
		if err != nil {
			t.Error(err)
		}
		if commands_output[i] != "" && response == nil {
			t.Fatal("Did not receive a response")
		} else if response.Logs != commands_output[i] {
			t.Fatalf("Did not receive correct response from command %s. Expected: %s Got: %s", commands[i], commands_output[i], response.Logs)
		}
	}
}

// Tests to ensure the agent can record properly
// Tests premature start and stop and proper recording
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
	t.Run("PreStop", TestPrematureStopRecording)
	if !t.Failed() {
		num_passed++
	}
	t.Run("Start", TestStartRecording)
	if t.Failed() {
		t.FailNow()
	} else {
		num_passed++
	}
	t.Run("ExtraStart", TestExtraStartRecording)
	if !t.Failed() {
		num_passed++
	}
	t.Run("Commands", TestCommands)
	if !t.Failed() {
		num_passed++
	}
	t.Run("Stop", TestStopRecording)
	if !t.Failed() {
		num_passed++
	}
}

// Tests attempting to stop recording before a recording has started
// The agent should prevent this from happening with an exception
// Should be run before agent.StartRecording
func TestPrematureStopRecording(t *testing.T) {
	num_tests++
	_, err := agent.StopRecording(ctx)
	if err == nil {
		t.Fatal("Did not prevent stopping a non-existant recording")
	} else {
		err_num := getError(err)
		err_expected := NOT_RECORDING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Test that a recording can be started without error
func TestStartRecording(t *testing.T) {
	num_tests++
	if err := agent.StartRecording(ctx, recording_name); err != nil {
		t.Error(err)
	}
}

// Tests attempting to start recording while a recording is in progress
// The agent should prevent this from happening with an exception
// Should be run after agent.StartRecording
func TestExtraStartRecording(t *testing.T) {
	num_tests++
	err := agent.StartRecording(ctx, recording_name)
	if err == nil {
		t.Fatal("Did not prevent starting a second concurrent recording")
	} else {
		err_num := getError(err)
		err_expected := RECORDING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Tests that a recording can be stopped without error
// Checks the returned recording name
func TestStopRecording(t *testing.T) {
	num_tests++
	recording, err := agent.StopRecording(ctx)
	if err != nil {
		t.Error(err)
	}
	if recording != nil {
		if recording.RecordingName != recording_name {
			t.Errorf("Did not return correct recording name. Expected: '%s' Got: '%s'", recording_name, recording.RecordingName)
		}
		snapshotName := fmt.Sprintf("%s/%s-rr-snp", recording.Location, recording_name)
		if recording.GetSnapshotFileName() != snapshotName {
			t.Errorf("Did not return correct snaphot name. Expected: '%s' Got: '%s'", snapshotName, recording.GetSnapshotFileName())
		}
		ndLogName := fmt.Sprintf("%s/%s-rr-nondet.log", recording.Location, recording_name)
		if recording.GetNdlogFileName() != ndLogName {
			t.Errorf("Did not return correct nondet log name. Expected: '%s' Got: '%s'", ndLogName, recording.GetNdlogFileName())
		}
	} else {
		t.Fatal("Did not return recording")
	}
}

var replay_agent controller.PandaReplayAgent

// Tests to ensure the agent can replay properly
// Tests premature stop and proper replay
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

	t.Run("PreStop", TestPrematureReplayStop)
	if !t.Failed() {
		num_passed++
	}
	t.Run("WrongReplay", TestNonexistantReplay)
	if !t.Failed() {
		num_passed++
	}
	t.Run("RunReplay", TestRunReplay)
	if !t.Failed() {
		num_passed++
	}
	t.Run("RunExtraReplay", TestRunExtraReplay)
	if !t.Failed() {
		num_passed++
	}
}

// Tests attempting to stop a replay when one is not in progress
// The agent should prevent this from happening with an exception
// Should be run before replay_agent.StartReplayAgent
func TestPrematureReplayStop(t *testing.T) {
	num_tests++
	_, err := replay_agent.StopReplay(ctx)
	if err == nil {
		t.Error("Did not prevent premature stop")
	} else {
		err_num := getError(err)
		err_expected := NOT_REPLAYING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Tests that a recording can be replayed without error
func TestRunReplay(t *testing.T) {
	num_tests++
	replay, err := replay_agent.StartReplayAgent(ctx, recording_name)
	if err != nil {
		t.Fatal(err)
	}
	if replay == nil {
		t.Fatal("Replay did not return")
	}
	if replay.Serial == "" {
		t.Error("Replay did not return serial")
	} else {
		// Check serial I/O, currently only works with single-line I/O
		serial := strings.Split(replay.Serial, "\r\n")
		for i, cmd := range serial {
			index := i / 2 // Lines alternate between command and response
			if i%2 == 0 {
				// Test for the serial command
				if !strings.HasSuffix(cmd, commands[index]) {
					t.Errorf("Expected: '%s' Got: '%s'", commands[index], cmd)
				}
			} else {
				// Test for response of serial command
				if !strings.HasSuffix(cmd, commands_output[index]) {
					t.Errorf("Expected: '%s' Got: '%s'", commands_output[index], cmd)
				}
			}
		}

	}
	if replay.Replay == "" {
		t.Error("Replay did not return execution")
	} else {
		// Check replay execution
		if !strings.Contains(replay.Replay, "Replay completed successfully") {
			t.Fatal("Replay did not complete successfully")
		}
	}
}

// Tests attempting to start a replay after one has started
// The agent should prevent this from happening with an exception
// Should be run before replay_agent.StartReplayAgent
func TestRunExtraReplay(t *testing.T) {
	num_tests++
	_, err := replay_agent.StartReplayAgent(ctx, recording_name)
	if err == nil {
		t.Fatal("Did not prevent extra replay")
	} else {
		err_num := getError(err)
		err_expected := RUNNING
		if err_num != err_expected {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
	}
}

// Tests attempting to start a replay that does not exist
// This should be prevented with an exception
func TestNonexistantReplay(t *testing.T) {
	num_tests++
	_, err := replay_agent.StartReplayAgent(ctx, " ")
	if err == nil {
		t.Fatal("Did not prevent nonexistant replay")
	} else if !strings.Contains(err.Error(), "Error in copying snapshot for replay") {
		// Error happens before agent enumeration
		t.Error("Incorrect error message")
	}
}
