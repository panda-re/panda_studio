package panda_controller

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type dockerGrpcPandaAgent struct {
	grpcAgent   PandaAgent
	cli         *docker.Client
	containerId *string
	sharedDir   *string
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
	sharedDir, err := os.MkdirTemp("/tmp/panda-studio", "panda-agent")
	if err != nil {
		return nil, err
	}

	agent := &dockerGrpcPandaAgent{
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

// StartRecording implements PandaAgent
func (pa *dockerGrpcPandaAgent) StartRecording(ctx context.Context, recordingName string) error {
	// Tell the agent to save the snapshot file in the shared folder rather than the current directory
	newRecordingName := fmt.Sprintf("./shared/%s", recordingName)
	return pa.grpcAgent.StartRecording(ctx, newRecordingName)
}

// StopRecording implements PandaAgent
func (pa *dockerGrpcPandaAgent) StopRecording(ctx context.Context) (*PandaAgentRecording, error) {
	recording, err := pa.grpcAgent.StopRecording(ctx)
	if err != nil {
		return nil, err
	}

	// Strip off the shared folder prefix so that paths are correct on this system
	recordingName := strings.Replace(recording.RecordingName, "./shared/", "", 1)

	new_recording := PandaAgentRecording{
		RecordingName: recordingName,
		Location:      *pa.sharedDir,
	}

	// Copy given image into shared directory
	//TODO: Replace destination with filesystem target
	nBytes, err := copyFileHelper(new_recording.GetSnapshotFileName(), "/tmp/panda-studio/", fmt.Sprintf("%s-rr-snp", new_recording.RecordingName))
	if (err != nil && err != io.EOF) || nBytes == 0 {
		return nil, errors.Wrap(err, "Error in copying snapshot")
	}

	nBytes, err = copyFileHelper(new_recording.GetNdlogFileName(), "/tmp/panda-studio/", fmt.Sprintf("%s-rr-nondet.log", new_recording.RecordingName))
	if (err != nil && err != io.EOF) || nBytes == 0 {
		return nil, errors.Wrap(err, "Error in copying nondeterministic log")
	}
	return &new_recording, nil
}

func (pa *dockerGrpcPandaAgent) StartReplay(ctx context.Context, recordingName string) (*PandaAgentRunCommandResult, error) {
	// TODO
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
	return pa.grpcAgent.StartReplay(ctx, recordingLocation)
}

func (pa *dockerGrpcPandaAgent) StopReplay(ctx context.Context) error {
	print("Stopping replay TODO")
	// TODO
	return pa.grpcAgent.StopReplay(ctx)
}

func (pa *dockerGrpcPandaAgent) startContainer(ctx context.Context) error {
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

func copyFileHelper(source_path string, destination_path string, file_name string) (int64, error) {

	sourceFileStat, err := os.Stat(source_path)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", source_path)
	}

	source, err := os.Open(source_path)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(destination_path + file_name)
	if err != nil {
		return 0, err
	}

	nBytes, err := io.Copy(destination, source)

	destination.Close()

	return nBytes, err
}
