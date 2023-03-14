package panda_controller

import (
	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/db/repos"
	"github.com/pkg/errors"
)

type PandaProgramExecutor struct {
	imageRepo repos.ImageRepository
}

type PandaProgramExecutorJob struct {
	agent PandaAgent
	image *models.Image
	prgram *models.InteractionProgram
}

type PandaProgramExecutorOptions struct {
	Image *models.Image
	Program *models.InteractionProgram
}

func (p *PandaProgramExecutor) NewExecutorJob(opts PandaProgramExecutorOptions) (*PandaProgramExecutorJob, error) {
	// Basic rundown of what will happen:
	// 1. pull information from the database
	//    - making this the caller's responsibility
	//	  - passed in as 'opts'
	// 2. open a stream to the file in blob storage
	// 3. create a panda instance using that file
	// 4. start the agent
	// 5. send the commands to the agent
	//    - offer an interface for real-time feedback, even if we don't currently use it
	//    - keep track of any recording files that are created
	// 6. stop the agent
	// 7. upload the recording files to blob storage
	return nil, errors.New("not implemented")
}