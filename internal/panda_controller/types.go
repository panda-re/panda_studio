package panda_controller

import (
	"context"
	"fmt"
)

type PandaAgent interface {
	// todo: add options such as architecture, image, networking, etc. to pass
	// to the agent
	StartAgent(ctx context.Context) error
	StopAgent(ctx context.Context) error
	RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error)
	StartRecording(ctx context.Context, recordingName string) error
	StopRecording(ctx context.Context) (*PandaAgentRecording, error)
	Close() error
}

type PandaReplayAgent interface {
	StartAgent(ctx context.Context) error
	StopAgent(ctx context.Context) error
	StartReplay(ctx context.Context, recordingName string) (*PandaAgentRunCommandResult, error)
	StopReplay(ctx context.Context) error
	Close() error
}

type PandaAgentRunCommandResult struct {
	Logs string
}

type PandaAgentRecording struct {
	RecordingName string
	Location      string
}

func (r *PandaAgentRecording) GetSnapshotFileName() string {
	return fmt.Sprintf("%s/%s-rr-snp", r.Location, r.RecordingName)
}

func (r *PandaAgentRecording) GetNdlogFileName() string {
	return fmt.Sprintf("%s/%s-rr-nondet.log", r.Location, r.RecordingName)
}
