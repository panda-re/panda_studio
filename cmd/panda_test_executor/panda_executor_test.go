package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	controller "github.com/panda-re/panda_studio/internal/panda_controller"
	"github.com/pkg/errors"
)

const QCOW_NAME = "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2"

var QCOW_LOCAL = fmt.Sprintf("/root/.panda/%s", QCOW_NAME)

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
	// Upon error, the following regex will be in the message
	re := regexp.MustCompile(`<ErrorCode\.\w+: `)
	// Regex splits right before error number
	nums := re.Split(err.Error(), -1)
	// Check if error matches regex. If not, it's not from agent
	if len(nums) < 2 {
		return -1
	}
	// Extract the error number from the message
	num, err := strconv.Atoi(nums[1][0:1])
	if err != nil {
		return -1
	}
	return num
}

// Checks if the error matches what is expected
// Prints a message if not
func checkError(err error, err_expected int, t *testing.T) {
	err_num := getError(err)
	if err_num != err_expected {
		if err_num != -1 {
			t.Errorf("Received wrong error. Expected: %s Got: %s", error_to_string[err_expected], error_to_string[err_num])
		}
		t.Error(err)
	}
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
	t.Run("Recording", TestRecord)
	t.Run("Replay", TestReplay)
}

var agent *controller.DockerPandaAgent

// Recording name for testing record and replay
const RECORDING_NAME string = "panda_executor_test"

const DEFAULT_QCOW_SIZE = 17711104

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
	agent, err = setupContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

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

	t.Run("PreCommand", TestPrematureCommand)
	// TODO absorb fail checks into tests
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
	// TODO test starting replay after agent start
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
		checkError(err, NOT_RUNNING, t)
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
			t.Error(err)
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
		checkError(err, RUNNING, t)
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
	agent, err = setupContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

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

	err = agent.StartAgent(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("PreStop", TestPrematureStopRecording)
	if !t.Failed() {
		num_passed++
	}
	t.Run("StartRecording", TestStartRecording)
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
	t.Run("StopRecording", TestStopRecording)
	if !t.Failed() {
		num_passed++
	}
	err = agent.StartRecording(ctx, "_")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("HangingStop", TestHangingRecording)
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
		checkError(err, NOT_RECORDING, t)
	}
}

// Test that a recording can be started without error
func TestStartRecording(t *testing.T) {
	num_tests++
	if err := agent.StartRecording(ctx, RECORDING_NAME); err != nil {
		t.Error(err)
	}
}

// Tests attempting to start recording while a recording is in progress
// The agent should prevent this from happening with an exception
// Should be run after agent.StartRecording
func TestExtraStartRecording(t *testing.T) {
	num_tests++
	err := agent.StartRecording(ctx, RECORDING_NAME)
	if err == nil {
		t.Fatal("Did not prevent starting a second concurrent recording")
	} else {
		checkError(err, RECORDING, t)
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
		if recording.Name() != RECORDING_NAME {
			t.Errorf("Did not return correct recording name. Expected: '%s' Got: '%s'", RECORDING_NAME, recording.Name())
		}
		ndl_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, recording.NdlogFilename())
		err = copyFileFromContainerHelper(ctx, recording.NdlogFilename(), ndl_dest, agent)
		if err != nil {
			panic(err)
		}
		snp_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, recording.SnapshotFilename())
		err = copyFileFromContainerHelper(ctx, recording.SnapshotFilename(), snp_dest, agent)
		if err != nil {
			panic(err)
		}
	} else {
		t.Fatal("Did not return recording")
	}
}

// Tests that a recording is stopped when agent is stopped
// The agent should stop the recording and raise a warning
// Should be run after agent.StartRecording
func TestHangingRecording(t *testing.T) {
	num_tests++
	err := agent.StopAgent(ctx)
	if err == nil {
		t.Fatal("Did not receive warning for stopping a hanging recording")
	} else {
		checkError(err, RECORDING, t)
	}
}

// Tests to ensure the agent can replay properly
// Tests premature stop and proper replay
func TestReplay(t *testing.T) {
	var err error
	agent, err = setupContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
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

	// Copy snapshot and nondet log to container for replay
	snp_name := fmt.Sprintf("%s-rr-snp", RECORDING_NAME)
	snp_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, snp_name)
	err = copyFileToContainerHelper(ctx, snp_dest, snp_name, agent)
	if err != nil {
		panic(err)
	}
	ndl_name := fmt.Sprintf("%s-rr-nondet.log", RECORDING_NAME)
	ndl_dest := fmt.Sprintf("%s/%s", controller.PANDA_STUDIO_TEMP_DIR, ndl_name)
	err = copyFileToContainerHelper(ctx, ndl_dest, ndl_name, agent)
	if err != nil {
		panic(err)
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
	_, err := agent.StopReplay(ctx)
	if err == nil {
		t.Error("Did not prevent premature stop")
	} else {
		checkError(err, NOT_REPLAYING, t)
	}
}

// Tests that a recording can be replayed without error
func TestRunReplay(t *testing.T) {
	num_tests++
	replay, err := agent.StartReplay(ctx, RECORDING_NAME)
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
	_, err := agent.StartReplay(ctx, RECORDING_NAME)
	if err == nil {
		t.Fatal("Did not prevent extra replay")
	} else {
		checkError(err, RUNNING, t)
	}
}

// Tests attempting to start a replay that does not exist
// This should be prevented with an exception
func TestNonexistantReplay(t *testing.T) {
	num_tests++
	_, err := agent.StartReplay(ctx, " ")
	if err == nil {
		t.Fatal("Did not prevent nonexistant replay")
	} else {
		checkError(err, REPLAYING, t)
	}
}

// Sets up the container for each test
// Uses the default config and qcow
func setupContainer(ctx context.Context) (*controller.DockerPandaAgent, error) {
	agent, err := controller.CreateDockerPandaAgent2(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create panda agent")
	}
	err = agent.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to agent")
	}
	err = copyFileToContainerHelper(ctx, QCOW_LOCAL, QCOW_NAME, agent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to copy qcow to agent")
	}
	return agent, nil
}

// ctx - context
// srcFilePath - file path on local machine
// dstFilePath - name of the file in the container
// agent - PandaAgent to container to copy into
func copyFileToContainerHelper(ctx context.Context, srcFilePath string, dstFilePath string, agent *controller.DockerPandaAgent) error {
	fileReader, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return err
	}
	err = agent.CopyFileToContainer(ctx, fileReader, fileInfo.Size(), dstFilePath)
	return err
}

// ctx - context
// srcFilePath - file path in container to copy from
// dstFilePath - file path on local machine to copy to
// agent - PandaAgent to container to copy from
func copyFileFromContainerHelper(ctx context.Context, srcFilePath string, dstFilePath string, agent *controller.DockerPandaAgent) error {
	src, err := agent.CopyFileFromContainer(ctx, srcFilePath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	nBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	if nBytes == 0 {
		return errors.New("did not copy file")
	}
	return nil
}
