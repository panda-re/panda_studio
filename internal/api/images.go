package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

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

//newName string, oldName string, dockerHubImageName string, size int
func (s *PandaStudioServer) CreateDerivedImage(ctx *gin.Context, imageId string, fileId string) {
	var deriveReq DeriveImageFileRequest
	err := ctx.BindJSON(&deriveReq)
	if err != nil {
		ctx.Error(errors.Wrap(err, "invalid request"))
		return
	}

	fileReader, err := s.imageRepo.OpenImageFile(ctx, db.ParseObjectID(imageId), db.ParseObjectID(fileId))
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer fileReader.Close()

	sharedDir, err := os.MkdirTemp("/tmp/panda-studio", "derive-image-tmp")
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	//create new file in shared dir to copy to
	destImageInSharedDir, err := os.Create(sharedDir + "/" + deriveReq.OldName)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer destImageInSharedDir.Close()

	//TODO: fix this, first arg needs to be a Write object
	nBytes, err := io.Copy(destImageInSharedDir, fileReader)
	if err != nil || nBytes == 0 {
		ctx.Error(errors.WithStack(err))
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd: []string{"docker build", "Dockerfile.derive-image",
			"--build-arg new_image=" + deriveReq.NewName + " ",
			"--build-arg base_image=" + deriveReq.OldName + " ",
			"--build-arg docker_image=" + deriveReq.DockerHubImageName + " ",
			"--build-arg size=" + deriveReq.Size},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	//retrieve derived image
	newImageFile, err := os.Open(sharedDir + "/" + deriveReq.newName)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}
	defer fileReader.Close()

	//upload image to repo
	created, err := s.imageRepo.Create(ctx, &models.Image{
		Name:        deriveReq.NewName,
		Description: "Derived from " + deriveReq.OldName,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	imageFile, err := s.imageRepo.CreateImageFile(ctx, &models.ImageFileCreateRequest{
		ImageID:  created.ID,
		FileName: created.Name,
		FileType: "qcow2",
	})
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	fileObj, err := s.imageRepo.UploadImageFile(ctx, &models.ImageFileUploadRequest{
		ImageId: created.ID,
		FileId:  imageFile.ID,
	}, newImageFile)
	if err != nil || fileObj == nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	return
}
