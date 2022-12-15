package panda_controller

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

type dockerGrpcPandaAgent struct {
	grpcAgent *grpcPandaAgent
	cli       *docker.Client
	containerId *string
}

const DOCKER_IMAGE = "pandare/pandaagent"

func CreateDefaultDockerPandaAgent(ctx context.Context) (PandaAgent, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	agent := &dockerGrpcPandaAgent{
		grpcAgent: nil,
		cli: cli,
		containerId: nil,
	}

	err = agent.startContainer(ctx)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

// Close implements PandaAgent
func (pa *dockerGrpcPandaAgent) Close() error {
	// Close grpc connection
	err := pa.grpcAgent.Close()
	if err != nil {
		return err
	}

	// Then stop/remove container
	err = pa.stopContainer(context.Background())
	if err != nil {
		return err
	}

	// Finally
	err = pa.cli.Close()
	if err != nil {
		return err
	}

	return nil
}

// RunCommand implements PandaAgent
func (pa *dockerGrpcPandaAgent) RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error) {
	return pa.grpcAgent.RunCommand(ctx, cmd)
}

// StartAgent implements PandaAgent
func (pa *dockerGrpcPandaAgent) StartAgent(ctx context.Context) error {
	return pa.grpcAgent.StartAgent(ctx)
}

// StopAgent implements PandaAgent
func (pa *dockerGrpcPandaAgent) StopAgent(ctx context.Context) error {
	return pa.grpcAgent.StopAgent(ctx)
}

func (pa *dockerGrpcPandaAgent) startContainer(ctx context.Context) error {
	// Create the container and save the name
	ccResp, err := pa.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty: true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{},
		// make sure the container is removed on exit
		AutoRemove: true,
	}, &network.NetworkingConfig{}, nil, "")

	if err != nil {
		return err
	}
	pa.containerId = &ccResp.ID

	// Start the container

	// Use ContainerAttach to get container logs
	// Use CopyFromContainer and CopyToContainer to copy files

	return nil
}

func (pa *dockerGrpcPandaAgent) stopContainer(ctx context.Context) error {
	// if container is not running, our job is done
	if pa.containerId == nil {
		return nil
	}

	err := pa.cli.ContainerRemove(ctx, *pa.containerId, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	return nil
}
