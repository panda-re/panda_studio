// - FindOne(context, ImageId)
//		- imageID should be accessible from image we are deriving from, passed down from front end
// - OpenImageFile(context, imageId, fileId)
//		- fileId comes from ?
//		- returns a io.ReadCloser
// - save the downloaded image file to a temp directory
//		- nBytes, err := io.Copy(destination, source)
// - share the temp directory with the docker container
// - Run the bash script on the image to create the new derived image
// - Pull the new image from the shared directory
// - Create(context, Image)
// - CreateImageFile(context, ImageFileCreateRequest) //steve running into an issue when uploading a qcow2 to the s3 bucket (resolved?)
// 		- ImageFileCreateRequest needs 
//			- ImageID
//			- file name (get from user input)
//			- filetype (get from ?)

//look at the branch steve is working on for read stuff (openimagefile example)

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

//takes in an existing image ID, a new name + size + docker image name
func (s *PandaStudioServer) CreateDerivedImage (ctx context.Context, imageId string, fileId string, newImageName string, newImageSize int) () {
	ctx := context.Background()
	//get the image from the repo
	image, err := s.imageRepo.FindOneImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectId(fileId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileReader, err := s.imageRepo.OpenImageFile(ctx, imageId, fileId)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer fileReader.Close()

	//create temp shared directory
	sharedDir, err := os.MkdirTemp("/tmp/panda-studio", "derive-image-tmp")
	if err != nil {
		return nil, err
	}

	//save image to temp shared directory
	nBytes, err := io.Copy(sharedDir, fileReader)
	if err != nil {
		return nil, err
	}

	//start docker container (Dockerfile.derive-image)
	//run bash script on image to create new derived image 
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: Ubuntu_2204, //TODO: get base ubuntu image for derive image job?
		AttachStdout: true,
		AttackStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: "bind",
				Source: sharedDir,
				Target: "panda/shared" //TODO: verify
			},
		},
		AutoRemove: true,
	}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return err
	}

	const containerId = &resp.ID

	err = cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	//once the job is finished, the new derived image will be in the shared directory

	//retrieve derived image from temp shared directory
	//sharedDir + newImageName
	contents, err := os.ReadFile(sharedDir + "/" + newImageName)
    if err != nil {
        fmt.Println("File reading error", err)
        return
    }

	//upload derived image to object storage
	fileObj, err = s.imageRepo.UploadImageFile(context, &models.ImageFileUploadRequest{
		ImageId: db.parseIbjectID(newImageId),
		FileId: fileObj.ID,
	}, contents) //TODO: make sure this is right
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

} 