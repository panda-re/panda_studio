package executor

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

const imageName = "pandare/panda"

func RunDocker() error {

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Pull the latest version of the image
	rc, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, rc)

	rc.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Volumes: make(map[string]struct{}),
		// Cmd: []string {"ls", "/"},
		Cmd: []string {"cat", "/readme.md"},
	}, &container.HostConfig{
		Mounts: []mount.Mount {
			{
				Type:          "bind",
				Source:        "/home/nick/panda_studio/readme.md",
				Target:        "/readme.md",
				ReadOnly:      true,
				Consistency:   "",
			},
		},
	}, &network.NetworkingConfig{}, nil, "hello-container")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}


	rc, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow: true,
	})
	if err != nil {
		return err
	}

	// output logging lines
	logReader := dlog.NewReader(rc)
	sc := bufio.NewScanner(logReader)

	for sc.Scan() {
		log.Println(sc.Text())
	}
	rc.Close()

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case body := <-statusCh:
		log.Printf("Container exited with code %d\n", body.StatusCode)
	}

	// Remove the container on exit
	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}