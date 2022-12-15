package panda_controller

import (
	"bufio"
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const imageName = "pandare/pandaagent"

//go:generate ./generate_stubs.sh

func RunDocker() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	/*
	// Pull the latest version of the image
	rc, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, rc)

	rc.Close()
	*/

	log.Println("Creating container...")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty: true,
		AttachStdout: true,
		AttachStderr: true,
		// Cmd: []string {"ls", "/"},
		// Cmd: []string {"find", "/panda_studio/executor/python"},
	}, &container.HostConfig{
		Mounts: []mount.Mount {
			{
				Type:          "bind",
				Source:        "/home/nick/panda_studio",
				Target:        "/panda_studio",
				ReadOnly:      true,
				Consistency:   "",
			},
		},
		PortBindings: nat.PortMap{
			"50051/tcp": []nat.PortBinding{{
				HostPort: "50051",
			}},
		},
	}, &network.NetworkingConfig{}, nil, "pandaagent-container")
	if err != nil {
		return err
	}

	log.Println("Starting container...")
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}


	go func() {
		rc, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow: true,
		})
		if err != nil {
			panic(err)
		}

		log.Println("getting logs")
		// output logging lines
		//logReader := dlog.NewReader(rc)
		sc := bufio.NewScanner(rc)

		for sc.Scan() {
			log.Printf("docker - %s\n", sc.Text())
		}
	}()

	time.Sleep(time.Second)
	// statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		return err
	// 	}
	// case body := <-statusCh:
	// 	log.Printf("Container exited with code %d\n", body.StatusCode)
	// }
	
	log.Println("Stopping container...")
	//err = cli.ContainerStop(ctx, resp.ID, nil)
	err = cli.ContainerStop(ctx, resp.ID, nil)
	if err != nil {
		return err
	}

	log.Println("Removing container...")
	// Remove the container on exit
	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	return nil
}
