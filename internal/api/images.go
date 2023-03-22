package api

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/panda-re/panda_studio/internal/db"
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/pkg/errors"
)

func (s *PandaStudioServer) FindAllImages(ctx *gin.Context) {
	// todo: allow search criteria
	images, err := s.imageRepo.FindAll(ctx)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, images)
}

func (s *PandaStudioServer) FindImageById(ctx *gin.Context, imageId string) {
	image, err := s.imageRepo.FindOne(ctx, db.ParseObjectID(imageId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, image)
}

func (s *PandaStudioServer) CreateImage(ctx *gin.Context) {
	var createReq CreateImageRequest
	err := ctx.BindJSON(&createReq)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	created, err := s.imageRepo.Create(ctx, &models.Image{
		Name:        *createReq.Name,
		Description: *createReq.Description,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (s *PandaStudioServer) DeleteImageById(ctx *gin.Context, imageId string) {
	deleted, err := s.imageRepo.DeleteOne(ctx, db.ParseObjectID(imageId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}

func (s *PandaStudioServer) UpdateImage(ctx *gin.Context, imageId string) {
	var updateReq CreateImageRequest
	err := ctx.BindJSON(&updateReq)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	updated, err := s.imageRepo.Update(ctx, &models.Image{
		Name:        *updateReq.Name,
		Description: *updateReq.Description,
		ID:          db.ParseObjectID(imageId),
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, updated)
}

func (s *PandaStudioServer) CreateImageFile(ctx *gin.Context, imageId string) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileObj, err := s.imageRepo.CreateImageFile(ctx, &models.ImageFileCreateRequest{
		ImageID:  db.ParseObjectID(imageId),
		FileName: form.Value["file_name"][0],
		FileType: form.Value["file_type"][0],
	})
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileReader, err := formFile.Open()
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer fileReader.Close()

	fileObj, err = s.imageRepo.UploadImageFile(ctx, &models.ImageFileUploadRequest{
		ImageId: db.ParseObjectID(imageId),
		FileId:  fileObj.ID,
	}, fileReader)

	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	// todo: convert to dto
	ctx.JSON(http.StatusOK, fileObj)
}

func (s *PandaStudioServer) DownloadImageFile(ctx *gin.Context, imageId ImageId, fileId FileId) {
	imgFile, err := s.imageRepo.FindOneImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}

	fileReader, err := s.imageRepo.OpenImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}
	defer fileReader.Close()

	ctx.DataFromReader(http.StatusOK, imgFile.Size, "application/octet-stream", fileReader, map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename=%s`, imgFile.FileName),
	})
}

func (s *PandaStudioServer) DeleteImageFile(ctx *gin.Context, imageId ImageId, fileId FileId) {
	imgFile, err := s.imageRepo.DeleteImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, imgFile)
}

func (s *PandaStudioServer) CreateDerivedImage(ctx *gin.Context, imageId string, fileId string, newName string, oldName string, dockerHubImageName string, size int) error {
	//get the image from the repo
	// image, err := s.imageRepo.FindOneImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	// if err != nil {
	// 	ctx.Error(errors.WithStack(err))
	// 	return err
	// }

	fileReader, err := s.imageRepo.OpenImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return err
	}
	defer fileReader.Close()

	//create temp shared directory
	sharedDir, err := os.MkdirTemp("/tmp/panda-studio", "derive-image-tmp")
	if err != nil {
		return err
	}

	//save image to temp shared directory
	nBytes, err := io.Copy(sharedDir, fileReader)
	if err != nil {
		return err
	}

	//run docker container for derive image job
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"docker build", "Dockerfile.derive-image", 
						"--build-arg new_image="+newName+" ", 
						"--build-arg base_image="+oldName+" ", 
						"--build-arg docker_image="+dockerHubImageName},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	//once the job is finished, the new derived image will be in the shared directory

	//retrieve derived image from temp shared directory
	//sharedDir + newImageName
	contents, err := os.ReadFile(sharedDir + "/" + newImageName)
	if err != nil {
		fmt.Println("File reading error", err)
		return err
	}

	//upload derived image to object storage
	fileObj, err := s.imageRepo.UploadImageFile(ctx, &models.ImageFileUploadRequest{
		ImageId: db.ParseObjectID(newImageId),
		FileId:  fileObj,
	}, contents) //TODO: make sure this is right
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return err
	}

	return nil
}
