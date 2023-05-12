package panda_controller

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"path"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/panda-re/panda_studio/internal/util"
	"github.com/pkg/errors"
)

const DI_DOCKER_IMAGE = "pandare/panda_di_agent"

type dockerPandaDIAgent struct {
	cli         *docker.Client
	containerId *string
	params      *PandaDIAgentParams
}
type PandaDIAgentParams struct {
	base_image_name string
	new_image       string
	resize          string //size to expand vm disk - "XGB" or "XMB", where X is the amount of space to add
	docker_image    string
	base_image      io.ReadCloser
	base_image_size int64
}

const CONTAINER_DI_DATA_DIR = "/root"

func CreateDockerPandaDIAgent(ctx context.Context, args PandaDIAgentParams) (*dockerPandaDIAgent, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	diAgent := &dockerPandaDIAgent{
		cli:         cli,
		containerId: nil,
		params:      &args,
	}

	return diAgent, nil
}

func (pdia *dockerPandaDIAgent) Run(ctx context.Context) error {
	fmt.Printf("Running dockerPandaDIAgent with params: %+v\n", pdia.params)
	err := pdia.createContainer(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Copying base image to container\n")
	err = pdia.copyImageFile(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Starting container\n")
	//start container
	err = pdia.startContainer(ctx)
	if err != nil {
		return err
	}

	//wait for container to stop
	wait, waitErr := pdia.cli.ContainerWait(ctx, *pdia.containerId, container.WaitConditionNextExit)
	select {
	case waitResult := <-wait:
		if waitResult.StatusCode != 0 {
			errDesc, err := pdia.getLogs(ctx)
			if err != nil {
				return errors.Wrap(err, "Could not get logs from container")
			}
			return errors.Errorf("Could not derive image:\n%s", errDesc)
		}
	case err = <-waitErr:
		//error channel outputted first
		return errors.Wrap(err, "Could not wait on container")
	}

	//copying file out, and destroying the container will be done by the caller
	return nil
}

func (pdia *dockerPandaDIAgent) CopyFileToContainer(ctx context.Context, data io.Reader, size int64, file string) error {
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
	return pdia.cli.CopyToContainer(ctx, *pdia.containerId, CONTAINER_DI_DATA_DIR, r, dockerTypes.CopyToContainerOptions{})
}

func (pdia *dockerPandaDIAgent) CopyFileFromContainer(ctx context.Context, file string) (io.ReadCloser, error) {
	r, w := io.Pipe()

	srcFile := path.Join(CONTAINER_DI_DATA_DIR, file)
	reader, pathStat, err := pdia.cli.CopyFromContainer(ctx, *pdia.containerId, srcFile)
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

func (pdia *dockerPandaDIAgent) OpenImageFile(ctx context.Context) (io.ReadCloser, error) {
	return pdia.CopyFileFromContainer(ctx, pdia.params.new_image)
}

func (pdia *dockerPandaDIAgent) Close() error {
	ctx := context.Background()

	pdia.params.base_image.Close()

	err := pdia.stopContainer(ctx)
	if err != nil {
		return err
	}

	err = pdia.cli.Close()
	if err != nil {
		return err
	}

	return nil
}

func (pdia *dockerPandaDIAgent) createContainer(ctx context.Context) error {
	if pdia.containerId != nil {
		return errors.New("derive image container already started")
	}

	ccResp, err := pdia.cli.ContainerCreate(ctx, &container.Config{
		Image:        DI_DOCKER_IMAGE,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Env: []string{
			fmt.Sprintf("BASE_IMAGE=%s", pdia.params.base_image_name),
			fmt.Sprintf("NEW_IMAGE=%s", pdia.params.new_image),
			fmt.Sprintf("SIZE=%s", pdia.params.resize),
			fmt.Sprintf("DOCKER_IMAGE=%s", pdia.params.docker_image),
		},
	}, &container.HostConfig{
		// make sure the container is removed on exit
		AutoRemove: false,
		// To allow for passing through KVM
		Privileged: true,
	}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return err
	}

	pdia.containerId = &ccResp.ID

	return nil
}

func (pdia *dockerPandaDIAgent) copyImageFile(ctx context.Context) error {
	err := pdia.CopyFileToContainer(ctx, pdia.params.base_image, pdia.params.base_image_size, pdia.params.base_image_name)
	if err != nil {
		return err
	}

	return nil
}

func (pdia *dockerPandaDIAgent) startContainer(ctx context.Context) error {
	err := pdia.cli.ContainerStart(ctx, *pdia.containerId, dockerTypes.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (pdia *dockerPandaDIAgent) stopContainer(ctx context.Context) error {
	if pdia.containerId == nil {
		return errors.New("derive image container not started")
	}

	// get the logs
	logs, err := pdia.getLogs(ctx)
	if err != nil {
		return err
	}
	fmt.Println(logs)

	err = pdia.cli.ContainerStop(ctx, *pdia.containerId, nil)
	if err != nil {
		return err
	}

	err = pdia.cli.ContainerRemove(ctx, *pdia.containerId, dockerTypes.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	pdia.containerId = nil
	return nil
}

func (pdia *dockerPandaDIAgent) getLogs(ctx context.Context) (string, error) {
	if pdia.containerId == nil {
		return "", errors.New("derive image container not started")
	}

	logs, err := pdia.cli.ContainerLogs(ctx, *pdia.containerId, dockerTypes.ContainerLogsOptions{
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
