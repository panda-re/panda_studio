package panda_controller

import "context"

type PandaAgent interface {
	// todo: add options such as architecture, image, networking, etc. to pass
	// to the agent
	StartAgent(ctx context.Context) error
	StopAgent(ctx context.Context) error
	RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error)
	Close() error
}

type PandaAgentRunCommandResult struct {
	Logs string
}