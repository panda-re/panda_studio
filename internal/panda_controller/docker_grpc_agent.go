package panda_controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

type dockerGrpcPandaAgent struct {
	grpcAgent PandaAgent
	cli       *docker.Client
	containerId *string
	sharedDir *string
}

const DOCKER_IMAGE = "pandare/panda_agent"
const DOCKER_GRPC_SOCKET_PATTERN = "unix://%s/panda-agent.sock"

func CreateDefaultDockerPandaAgent(ctx context.Context) (PandaAgent, error) {
	// Connect to docker daemon
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Create a shared temporary directory
	sharedDir, err := os.MkdirTemp("", "panda-agent")
	if err != nil {
		return nil, err
	}

	agent := &dockerGrpcPandaAgent{
		grpcAgent: nil,
		cli: cli,
		containerId: nil,
		sharedDir: &sharedDir,
	}

	// Start the container
	err = agent.startContainer(ctx)
	if err != nil {
		return nil, err
	}

	// Wait for container startup
	time.Sleep(time.Millisecond*1000)

	// Connect to grpc over unix socket
	grpcSocket := fmt.Sprintf(DOCKER_GRPC_SOCKET_PATTERN, *agent.sharedDir)
	grpcAgent, err := CreateGrpcPandaAgent(grpcSocket)
	if err != nil {
		return nil, err
	}	

	agent.grpcAgent = grpcAgent

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

	// Close docker connection
	err = pa.cli.Close()
	if err != nil {
		return err
	}

	// Remove temp dir
	if pa.sharedDir != nil {
		err = os.RemoveAll(*pa.sharedDir)
		if err != nil {
			return err
		}
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
		Image: DOCKER_IMAGE,
		Tty: true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: "bind",
				Source: *pa.sharedDir,
				Target: "/panda/shared",
			},
			// So PANDA doesn't need to download the same image
			{
				Type: "bind",
				Source: "/home/nick/.panda",
				Target: "/root/.panda",
			},
		},
		// make sure the container is removed on exit
		AutoRemove: true,
	}, &network.NetworkingConfig{}, nil, "")

	if err != nil {
		return err
	}
	pa.containerId = &ccResp.ID

	// Start the container
	err = pa.cli.ContainerStart(ctx, *pa.containerId, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

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
