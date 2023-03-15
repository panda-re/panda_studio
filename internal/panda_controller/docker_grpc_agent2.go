package panda_controller

import (
	"context"

	docker "github.com/docker/docker/client"
)

type dockerGrpcPandaAgent2 struct {
	grpcAgent   PandaAgent
	cli         *docker.Client
	containerId *string
	sharedDir   *string
}

var _ PandaAgent = &dockerGrpcPandaAgent2{}

type StartPandaOpts struct {
	// In the future this should include stuff like architecture, image, etc.
	
}

func CreateDockerPandaAgent2(ctx context.Context) (*dockerGrpcPandaAgent2, error) {
	// Connect to docker daemon
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Initialize agent
	agent := &dockerGrpcPandaAgent2{
		grpcAgent:   nil,
		cli:         cli,
		containerId: nil,
		sharedDir:   nil,
	}

	return agent, nil
}

// Close implements PandaAgent
func (*dockerGrpcPandaAgent2) Close() error {
	panic("unimplemented")
}

// RunCommand implements PandaAgent
func (*dockerGrpcPandaAgent2) RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error) {
	panic("unimplemented")
}

// StartAgent implements PandaAgent
// Starts the agent using the default x86_64 image
func (*dockerGrpcPandaAgent2) StartAgent(ctx context.Context) error {
	panic("unimplemented")
}

func (*dockerGrpcPandaAgent2) StartAgentWithImage(ctx context.Context, opts *StartPandaOpts) error {
	panic("unimplemented")
}

// StartRecording implements PandaAgent
func (*dockerGrpcPandaAgent2) StartRecording(ctx context.Context, recordingName string) error {
	panic("unimplemented")
}

// StopAgent implements PandaAgent
func (*dockerGrpcPandaAgent2) StopAgent(ctx context.Context) error {
	panic("unimplemented")
}

// StopRecording implements PandaAgent
func (*dockerGrpcPandaAgent2) StopRecording(ctx context.Context) (*PandaAgentRecording, error) {
	panic("unimplemented")
}