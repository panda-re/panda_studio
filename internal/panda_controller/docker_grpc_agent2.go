package panda_controller

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/panda-re/panda_studio/internal/util"
	"github.com/panda-re/panda_studio/panda_agent/pb"
)

type DockerGrpcPandaAgent2 struct {
	grpcAgent   PandaAgent
	cli         *docker.Client
	containerId *string
}

var _ PandaAgent = &DockerGrpcPandaAgent2{}

const PANDA_STUDIO_TEMP_DIR = "/tmp/panda-studio"

func CreateDockerPandaAgent2(ctx context.Context) (*DockerGrpcPandaAgent2, error) {
	// Connect to docker daemon
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Initialize agent
	agent := &DockerGrpcPandaAgent2{
		grpcAgent:   nil,
		cli:         cli,
		containerId: nil,
	}

	return agent, nil
}

func (pa *DockerGrpcPandaAgent2) Connect(ctx context.Context) error {
	// start the container
	err := pa.startContainer(ctx)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	// connect to the grpc agent

	// ping the agent to ensure successful connection
	// panic("unimplemented")
	return nil
}

// See documentation for docker CopyToContainer
func (pa *DockerGrpcPandaAgent2) CopyFileToContainer(ctx context.Context, data io.Reader, size int64, file string) error {
	fmt.Println("Container ID: " + *pa.containerId)
	r, w := io.Pipe()
	tarWriter := tar.NewWriter(w)
	go func() {
		defer tarWriter.Close()
		defer w.Close()

		err := tarWriter.WriteHeader(&tar.Header{
			Name: file,
			Mode: 0644,
			Size: size,
		})
		if err != nil {
			return
		}
		_, err = util.Copy(ctx, tarWriter, data)
		if err != nil {
			return
		}
	}()
	return pa.cli.CopyToContainer(ctx, *pa.containerId, "/panda/shared", r, dockerTypes.CopyToContainerOptions{})
}

// See documentation for docker CopyFromContainer
func (pa *DockerGrpcPandaAgent2) CopyFileFromContainer(ctx context.Context, file string) (io.ReadCloser, error) {
	r, w := io.Pipe()

	srcFile := path.Join("/panda/shared", file)
	reader, pathStat, err := pa.cli.CopyFromContainer(ctx, *pa.containerId, srcFile)
	if err != nil {
		return nil, err
	}

	go func() {
		defer w.Close()
		defer reader.Close()

		tarReader := tar.NewReader(reader)
		header, err := tarReader.Next()
		if err != nil {
			return
		}
		if header.Name != pathStat.Name || pathStat.Mode.IsDir() {
			return
		}
		_, err = util.Copy(ctx, w, tarReader)
		if err != nil {
			return
		}
	}()


	return r, nil
}

// Close implements PandaAgent
func (pa *DockerGrpcPandaAgent2) Close() error {
	// Close grpc connection
	if pa.grpcAgent != nil {
		err := pa.grpcAgent.Close()
		if err != nil {
			return err
		}
	}

	// Stop and remove container
	err := pa.stopContainer(context.Background())
	if err != nil {
		return err
	}

	// Close docker client
	err = pa.cli.Close()
	if err != nil {
		return err
	}

	return nil
}

// RunCommand implements PandaAgent
func (*DockerGrpcPandaAgent2) RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error) {
	panic("unimplemented")
}

// StartAgent implements PandaAgent
// Starts the agent using the default x86_64 image
func (*DockerGrpcPandaAgent2) StartAgent(ctx context.Context) error {
	panic("unimplemented")
}

func (*DockerGrpcPandaAgent2) StartAgentWithOpts(ctx context.Context, opts *pb.StartAgentRequest) error {
	panic("unimplemented")
}

// StartRecording implements PandaAgent
func (*DockerGrpcPandaAgent2) StartRecording(ctx context.Context, recordingName string) error {
	panic("unimplemented")
}

// StopAgent implements PandaAgent
func (*DockerGrpcPandaAgent2) StopAgent(ctx context.Context) error {
	panic("unimplemented")
}

// StopRecording implements PandaAgent
func (*DockerGrpcPandaAgent2) StopRecording(ctx context.Context) (*PandaAgentRecording, error) {
	panic("unimplemented")
}

func (pa *DockerGrpcPandaAgent2) startContainer(ctx context.Context) error {
	if pa.containerId != nil {
		return errors.New("container already started")
	}
	// Create the container and save the name
	ccResp, err := pa.cli.ContainerCreate(ctx, &container.Config{
		Image:        DOCKER_IMAGE,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		// Instead of mounting we will copy files to the container
		/*
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: *pa.sharedDir,
				Target: "/panda/shared",
			},
		},
		*/
		// make sure the container is removed on exit
		// AutoRemove: true,
		AutoRemove: false,
	}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return err
	}

	pa.containerId = &ccResp.ID

	// Start the container
	err = pa.cli.ContainerStart(ctx, *pa.containerId, dockerTypes.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (pa *DockerGrpcPandaAgent2) stopContainer(ctx context.Context) error {
	if pa.containerId == nil {
		return errors.New("container not started")
	}
	err := pa.cli.ContainerStop(ctx, *pa.containerId, nil)
	if err != nil {
		return err
	}
	pa.containerId = nil
	return nil
}