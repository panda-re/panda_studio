package panda_controller

// ! Most of the code in this file is duplicated from docker_grpc_agent.go and
// ! should be refactored for code reuse

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type dockerGrpcPandaReplayAgent struct {
	grpcAgent   PandaReplayAgent
	cli         *docker.Client
	containerId *string
	sharedDir   *string
}

func CreateReplayDockerPandaAgent(ctx context.Context) (PandaReplayAgent, error) {
	// Connect to docker daemon
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Create a shared temporary directory
	sharedDir, err := os.MkdirTemp("/tmp/panda-studio", "panda-agent")
	if err != nil {
		return nil, err
	}

	agent := &dockerGrpcPandaReplayAgent{
		grpcAgent:   nil,
		cli:         cli,
		containerId: nil,
		sharedDir:   &sharedDir,
	}

	// Start the container
	err = agent.startContainer(ctx)
	if err != nil {
		return nil, err
	}

	// Wait for container startup - we need a better method
	time.Sleep(time.Millisecond * 2000)

	// Connect to grpc over unix socket
	grpcSocket := fmt.Sprintf(DOCKER_GRPC_SOCKET_PATTERN, *agent.sharedDir)
	grpcAgent, err := CreateGrpcPandaAgent(grpcSocket)
	if err != nil {
		return nil, err
	}

	agent.grpcAgent = grpcAgent.(PandaReplayAgent)

	return agent, nil
}

// Close implements PandaReplayAgent
func (pa *dockerGrpcPandaReplayAgent) Close() error {
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

// StopAgent implements PandaReplayAgent
func (pa *dockerGrpcPandaReplayAgent) StopAgent(ctx context.Context) error {
	return pa.grpcAgent.StopAgent(ctx)
}

// StartReplay implements PandaReplayAgent
func (pa *dockerGrpcPandaReplayAgent) StartReplayAgent(ctx context.Context, recordingName string) (*PandaAgentReplayResult, error) {
	// Copy file into shared directory
	sharedFolder := fmt.Sprintf("%s/", *pa.sharedDir)
	snapshotName := fmt.Sprintf("%s-rr-snp", recordingName)
	nBytes, err := copyFileHelper(fmt.Sprintf("/tmp/panda-studio/%s", snapshotName), sharedFolder, snapshotName)
	if (err != nil && err != io.EOF) || nBytes == 0 {
		return nil, errors.Wrap(err, "Error in copying snapshot for replay")
	}

	ndLogName := fmt.Sprintf("%s-rr-nondet.log", recordingName)
	nBytes, err = copyFileHelper(fmt.Sprintf("/tmp/panda-studio/%s", ndLogName), sharedFolder, ndLogName)
	if (err != nil && err != io.EOF) || nBytes == 0 {
		return nil, errors.Wrap(err, "Error in copying nondeterministic log for replay")
	}
	// Location of the recording for PANDA/container
	recordingLocation := fmt.Sprintf("./shared/%s", recordingName)
	return pa.grpcAgent.StartReplayAgent(ctx, recordingLocation)
}

// StopReplay implements PandaReplayAgent
func (pa *dockerGrpcPandaReplayAgent) StopReplay(ctx context.Context) (*PandaAgentReplayResult, error) {
	return pa.grpcAgent.StopReplay(ctx)
}

func (pa *dockerGrpcPandaReplayAgent) startContainer(ctx context.Context) error {
	// Create the container and save the name
	ccResp, err := pa.cli.ContainerCreate(ctx, &container.Config{
		Image:        DOCKER_IMAGE,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: *pa.sharedDir,
				Target: "/panda/shared",
			},
			// So PANDA doesn't need to download the same image
			{
				Type:   "bind",
				Source: "/root/.panda",
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

func (pa *dockerGrpcPandaReplayAgent) stopContainer(ctx context.Context) error {
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