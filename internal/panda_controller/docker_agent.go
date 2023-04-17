package panda_controller

import (
	"archive/tar"
	"context"
	"errors"
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
)

type dockerPandaAgent struct {
	grpcAgent   PandaAgent
	cli         *docker.Client
	containerId *string
	sharedDir   *string
}

var _ PandaAgent = &dockerPandaAgent{}

const PANDA_STUDIO_TEMP_DIR = "/tmp/panda-studio"
const CONTAINER_SHARED_DIR = "/panda/shared"
const CONTAINER_DATA_DIR = "/panda/data"
const DOCKER_IMAGE = "pandare/panda_agent"

func CreateDockerPandaAgent2(ctx context.Context) (*dockerPandaAgent, error) {
	// Connect to docker daemon
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Initialize agent
	agent := &dockerPandaAgent{
		grpcAgent:   nil,
		cli:         cli,
		containerId: nil,
	}

	return agent, nil
}

func (pa *dockerPandaAgent) Connect(ctx context.Context) error {
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

func (pa *dockerPandaAgent) connectGrpc(ctx context.Context) error {
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
func (pa *dockerPandaAgent) CopyFileToContainer(ctx context.Context, data io.Reader, size int64, file string) error {
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
func (pa *dockerPandaAgent) CopyFileFromContainer(ctx context.Context, file string) (io.ReadCloser, error) {
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

func (pa *dockerPandaAgent) OpenRecordingSnapshot(ctx context.Context, recordingName string) (io.ReadCloser, error) {
	snapshotFileName := fmt.Sprintf("%s-rr-snp", recordingName)
	return pa.CopyFileFromContainer(ctx, snapshotFileName)
}

func (pa *dockerPandaAgent) OpenRecordingNdlog(ctx context.Context, recordingName string) (io.ReadCloser, error) {
	ndlogFileName := fmt.Sprintf("%s-rr-nondet.log", recordingName)
	return pa.CopyFileFromContainer(ctx, ndlogFileName)
}

// Close implements PandaAgent
func (pa *dockerPandaAgent) Close() error {
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
func (pa *dockerPandaAgent) RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error) {
	return pa.grpcAgent.RunCommand(ctx, cmd)
}

// StartAgent implements PandaAgent
func (pa *dockerPandaAgent) StartAgent(ctx context.Context) error {
	config := pb.PandaConfig{
		QcowFileName: "bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2",
		Arch:         "x86_64",
		Os:           "linux-64-ubuntu:4.15.0-72-generic-noaslr-nokaslr",
		Prompt:       "root@ubuntu:.*# ",
		Cdrom:        "",
		Snapshot:     "root",
		Memory:       "1024",
		ExtraArgs:    "-display none",
	}
	return pa.grpcAgent.StartAgentWithOpts(ctx, &pb.StartAgentRequest{Config: &config})
}

func (pa *dockerPandaAgent) StartAgentWithOpts(ctx context.Context, opts *pb.StartAgentRequest) error {
	return pa.grpcAgent.StartAgentWithOpts(ctx, opts)
}

// StartRecording implements PandaAgent
func (pa *dockerPandaAgent) StartRecording(ctx context.Context, recordingName string) error {
	return pa.grpcAgent.StartRecording(ctx, recordingName)
}

// StopAgent implements PandaAgent
func (pa *dockerPandaAgent) StopAgent(ctx context.Context) error {
	return pa.grpcAgent.StopAgent(ctx)
}

// StopRecording implements PandaAgent
func (pa *dockerPandaAgent) StopRecording(ctx context.Context) (PandaAgentRecording, error) {
	innerRec, err := pa.grpcAgent.StopRecording(ctx)
	if err != nil {
		return nil, err
	}

	recording := innerRec.(*GenericPandaAgentRecordingConcrete)

	newRecording := DockerPandaAgentRecording{
		GenericPandaAgentRecordingConcrete: *recording,
		agent:                              pa,
	}

	return &newRecording, nil
}

func (pa *dockerPandaAgent) createTempDir() error {
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

func (pa *dockerPandaAgent) removeTempDir() error {
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

func (pa *dockerPandaAgent) startContainer(ctx context.Context) error {
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

func (pa *dockerPandaAgent) stopContainer(ctx context.Context) error {
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
