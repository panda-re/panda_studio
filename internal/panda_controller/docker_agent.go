package panda_controller

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/panda-re/panda_studio/internal/util"
	"github.com/panda-re/panda_studio/panda_agent/pb"
	"github.com/pkg/errors"
)

type DockerPandaAgent struct {
	grpcAgent   PandaAgent
	cli         *docker.Client
	containerId *string
	sharedDir   *string
}

var _ PandaAgent = &DockerPandaAgent{}

// See https://github.com/panda-re/panda/blob/dev/panda/python/core/pandare/qcows.json
var DEFAULT_CONFIG = pb.PandaConfig{
	Qcowfilename: "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2",
	Arch:         "x86_64",
	Os:           "linux-64-ubuntu:4.15.0-72-generic-noaslr-nokaslr",
	Prompt:       "root@ubuntu:.*#",
	Cdrom:        "ide1-cd0",
	Snapshot:     "root",
	Memory:       "1024",
	Extraargs:    "-display none",
}

const PANDA_STUDIO_TEMP_DIR = "/tmp/panda-studio"
const CONTAINER_SHARED_DIR = "/panda/shared"
const CONTAINER_DATA_DIR = "/panda/data"
const DOCKER_IMAGE = "pandare/panda_agent"

// Creates a PandaAgent inside a Docker client
func CreateDockerPandaAgent2(ctx context.Context) (*DockerPandaAgent, error) {
	// Connect to docker daemon
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Initialize agent
	agent := &DockerPandaAgent{
		grpcAgent:   nil,
		cli:         cli,
		containerId: nil,
	}

	return agent, nil
}

// Starts the Docker container and connects it to the gRPC agent
func (pa *DockerPandaAgent) Connect(ctx context.Context) error {
	// start the container
	err := pa.startContainer(ctx)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	// connect to the grpc agent
	err = pa.connectGrpc(ctx)
	if err != nil {
		return err
	}

	// ping the agent to ensure successful connection
	// panic("unimplemented")
	return nil
}

// Connects to the gRPC agent in the Docker container that is spun up
func (pa *DockerPandaAgent) connectGrpc(ctx context.Context) error {
	grpcSocketPath := path.Join(*pa.sharedDir, "panda-agent.sock")
	grpcSocketUrl := fmt.Sprintf("unix://%s", grpcSocketPath)
	grpcAgent, err := CreateGrpcPandaAgent(grpcSocketUrl)
	if err != nil {
		return err
	}
	pa.grpcAgent = grpcAgent.(PandaAgent)

	return nil
}

// See documentation for docker CopyToContainer
func (pa *DockerPandaAgent) CopyFileToContainer(ctx context.Context, data io.Reader, size int64, file string) error {
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
	return pa.cli.CopyToContainer(ctx, *pa.containerId, CONTAINER_DATA_DIR, r, dockerTypes.CopyToContainerOptions{})
}

// See documentation for docker CopyFromContainer
func (pa *DockerPandaAgent) CopyFileFromContainer(ctx context.Context, file string) (io.ReadCloser, error) {
	r, w := io.Pipe()

	srcFile := path.Join(CONTAINER_DATA_DIR, file)
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

func (pa *DockerPandaAgent) OpenRecordingSnapshot(ctx context.Context, recordingName string) (io.ReadCloser, error) {
	snapshotFileName := fmt.Sprintf("%s-rr-snp", recordingName)
	return pa.CopyFileFromContainer(ctx, snapshotFileName)
}

func (pa *DockerPandaAgent) OpenRecordingNdlog(ctx context.Context, recordingName string) (io.ReadCloser, error) {
	ndlogFileName := fmt.Sprintf("%s-rr-nondet.log", recordingName)
	return pa.CopyFileFromContainer(ctx, ndlogFileName)
}

// Close implements PandaAgent
func (pa *DockerPandaAgent) Close() error {
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

	// Ensure temp dir is removed
	err = pa.removeTempDir()
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
func (pa *DockerPandaAgent) RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error) {
	resp, err := pa.grpcAgent.RunCommand(ctx, cmd)
	if err != nil {
		logs, logErr := pa.getLogs(ctx)
		if logErr != nil {
			logs = fmt.Sprintf("Failed to get logs: %s", logErr.Error())
		}
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to run command %s. Logs: %s", cmd, logs))
	}

	return resp, nil
}

// StartAgent implements PandaAgent
// Uses x86_64 generic defaults
func (pa *DockerPandaAgent) StartAgent(ctx context.Context) error {
	return pa.grpcAgent.StartAgentWithOpts(ctx, &pb.StartAgentRequest{Config: &DEFAULT_CONFIG})
}

// StartAgentWithOpts implements PandaAgent
func (pa *DockerPandaAgent) StartAgentWithOpts(ctx context.Context, opts *pb.StartAgentRequest) error {
	err := pa.grpcAgent.StartAgentWithOpts(ctx, opts)
	if err != nil {
		logs, logErr := pa.getLogs(ctx)
		if logErr != nil {
			logs = fmt.Sprintf("Failed to get logs: %s", logErr.Error())
		}
		return errors.Wrap(err, fmt.Sprintf("Failed to start agent. Logs: %s", logs))
	}

	return nil
}

// StartRecording implements PandaAgent
func (pa *DockerPandaAgent) StartRecording(ctx context.Context, recordingName string) error {
	err := pa.grpcAgent.StartRecording(ctx, recordingName)
	if err != nil {
		logs, logErr := pa.getLogs(ctx)
		if logErr != nil {
			logs = fmt.Sprintf("Failed to get logs: %s", logErr.Error())
		}
		return errors.Wrap(err, fmt.Sprintf("Failed to start recording. Logs: %s", logs))
	}

	return nil
}

// StopAgent implements PandaAgent
func (pa *DockerPandaAgent) StopAgent(ctx context.Context) error {
	err := pa.grpcAgent.StopAgent(ctx)
	if err != nil {
		logs, logErr := pa.getLogs(ctx)
		if logErr != nil {
			logs = fmt.Sprintf("Failed to get logs: %s", logErr.Error())
		}
		return errors.Wrap(err, fmt.Sprintf("Failed to stop agent. Logs: %s", logs))
	}

	return nil
}

// StopRecording implements PandaAgent
func (pa *DockerPandaAgent) StopRecording(ctx context.Context) (PandaAgentRecording, error) {
	innerRec, err := pa.grpcAgent.StopRecording(ctx)
	if err != nil {
		logs, logErr := pa.getLogs(ctx)
		if logErr != nil {
			logs = fmt.Sprintf("Failed to get logs: %s", logErr.Error())
		}
		return nil, errors.Wrap(err, fmt.Sprintf("Stopping recording failed. Logs: %s", logs))
	}

	recording := innerRec.(*GenericPandaAgentRecordingConcrete)

	newRecording := DockerPandaAgentRecording{
		GenericPandaAgentRecordingConcrete: *recording,
		agent:                              pa,
	}

	return &newRecording, nil
}

// StartReplay implements PandaAgent
// Uses x86_64 generic defaults
func (pa *DockerPandaAgent) StartReplay(ctx context.Context, recordingName string) (*PandaAgentReplayResult, error) {
	return pa.grpcAgent.StartReplayWithOpts(ctx, &pb.StartAgentRequest{Config: &DEFAULT_CONFIG}, recordingName)
}

// StartReplayWithOpts implements PandaAgent
func (pa *DockerPandaAgent) StartReplayWithOpts(ctx context.Context, opts *pb.StartAgentRequest, recordingName string) (*PandaAgentReplayResult, error) {
	return pa.grpcAgent.StartReplayWithOpts(ctx, opts, recordingName)
}

// StopReplay implements PandaAgent
func (pa *DockerPandaAgent) StopReplay(ctx context.Context) (*PandaAgentReplayResult, error) {
	return pa.grpcAgent.StopReplay(ctx)
}

func (pa *DockerPandaAgent) createTempDir() error {
	if pa.sharedDir != nil {
		return nil
	}

	// Create temp dir
	sharedDir, err := os.MkdirTemp(PANDA_STUDIO_TEMP_DIR, "panda-studio-*")
	if err != nil {
		return err
	}

	pa.sharedDir = &sharedDir
	return nil
}

func (pa *DockerPandaAgent) removeTempDir() error {
	if pa.sharedDir == nil {
		return nil
	}

	// Remove temp dir
	err := os.RemoveAll(*pa.sharedDir)
	if err != nil {
		return err
	}

	pa.sharedDir = nil
	return nil
}

func (pa *DockerPandaAgent) startContainer(ctx context.Context) error {
	if pa.containerId != nil {
		return errors.New("container already started")
	}

	// Ensure temp dir exists
	err := pa.createTempDir()
	if err != nil {
		return err
	}

	// Create the container and save the name
	ccResp, err := pa.cli.ContainerCreate(ctx, &container.Config{
		Image:        DOCKER_IMAGE,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		// Instead of mounting we will copy files to the container
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: *pa.sharedDir,
				Target: CONTAINER_SHARED_DIR,
			},
		},
		// make sure the container is removed on exit
		AutoRemove: true,
		// AutoRemove: false,
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

func (pa *DockerPandaAgent) stopContainer(ctx context.Context) error {
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

func (pa *DockerPandaAgent) getLogs(ctx context.Context) (string, error) {
	if pa.containerId == nil {
		return "", errors.New("container not started")
	}

	logs, err := pa.cli.ContainerLogs(ctx, *pa.containerId, dockerTypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		return "", err
	}

	logsBytes, err := io.ReadAll(logs)
	if err != nil {
		return "", err
	}

	return string(logsBytes), nil
}
