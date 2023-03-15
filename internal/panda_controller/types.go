package panda_controller

import (
	"context"
	"fmt"

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
	StopRecording(ctx context.Context) (*PandaAgentRecording, error)
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

type PandaAgentRecording struct {
	RecordingName string
	Location      string
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

func (r *PandaAgentRecording) GetSnapshotFileName() string {
	return fmt.Sprintf("%s/%s-rr-snp", r.Location, r.RecordingName)
}

func (r *PandaAgentRecording) GetNdlogFileName() string {
	return fmt.Sprintf("%s/%s-rr-nondet.log", r.Location, r.RecordingName)
}
