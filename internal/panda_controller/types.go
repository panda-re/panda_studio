package panda_controller

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/panda-re/panda_studio/panda_agent/pb"
)

type StartAgentRequest = pb.StartAgentRequest

type PandaAgent interface {
	// todo: add options such as architecture, image, networking, etc. to pass
	// to the agent
	StartAgent(ctx context.Context) error
	StartAgentWithOpts(ctx context.Context, opts *StartAgentRequest) error
	StopAgent(ctx context.Context) error
	RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error)
	StartRecording(ctx context.Context, recordingName string) error
	StopRecording(ctx context.Context) (PandaAgentRecording, error)
	//SendNetworkCommand(ctx context.Context, network_request *NetworkRequest) (*NetworkResponse, error)
	Close() error
}

type PandaReplayAgent interface {
	StartReplayAgent(ctx context.Context, recordingName string) (*PandaAgentReplayResult, error)
	StopAgent(ctx context.Context) error
	StopReplay(ctx context.Context) (*PandaAgentReplayResult, error)
	Close() error
}

type PandaAgentRunCommandResult struct {
	Logs string
}

type PandaAgentReplayResult struct {
	Serial string // Captured serial through PANDA callback
	Replay string // Replay execution through redirected output to file
}

type NetworkRequest struct {
	SocketType   string
	Port         int32
	Application  string
	Command      string
	CustomPacket string
}

type NetworkResponse struct {
	StatusCode int32
	Output     string
}

type PandaAgentRecording interface {
	Name() string
	SnapshotFilename() string
	NdlogFilename() string
	OpenSnapshot(ctx context.Context) (io.ReadCloser, error)
	OpenNdlog(ctx context.Context) (io.ReadCloser, error)
}

type GenericPandaAgentRecordingConcrete struct {
	RecordingName string
}

var _ PandaAgentRecording = &GenericPandaAgentRecordingConcrete{}


func (r *GenericPandaAgentRecordingConcrete) Name() string {
	return r.RecordingName
}

func (r *GenericPandaAgentRecordingConcrete) SnapshotFilename() string {
	return fmt.Sprintf("%s-rr-snp", r.RecordingName)
}

func (r *GenericPandaAgentRecordingConcrete) NdlogFilename() string {
	return fmt.Sprintf("%s-rr-nondet.log", r.RecordingName)
}

func (r *GenericPandaAgentRecordingConcrete) OpenSnapshot(ctx context.Context) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (r *GenericPandaAgentRecordingConcrete) OpenNdlog(ctx context.Context) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

type DockerPandaAgentRecording struct {
	GenericPandaAgentRecordingConcrete
	agent *dockerPandaAgent
}

var _ PandaAgentRecording = &DockerPandaAgentRecording{}

func (r *DockerPandaAgentRecording) Name() string {
	return r.GenericPandaAgentRecordingConcrete.Name()
}

func (r *DockerPandaAgentRecording) SnapshotFilename() string {
	return r.GenericPandaAgentRecordingConcrete.SnapshotFilename()
}

func (r *DockerPandaAgentRecording) NdlogFilename() string {
	return r.GenericPandaAgentRecordingConcrete.NdlogFilename()
}

func (r *DockerPandaAgentRecording) OpenSnapshot(ctx context.Context) (io.ReadCloser, error) {
	return r.agent.CopyFileFromContainer(ctx, r.SnapshotFilename())
}

func (r *DockerPandaAgentRecording) OpenNdlog(ctx context.Context) (io.ReadCloser, error) {
	return r.agent.CopyFileFromContainer(ctx, r.NdlogFilename())
}