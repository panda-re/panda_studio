package panda_controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/db/repos"
	"github.com/pkg/errors"
)

type PandaProgramExecutor struct {
	imageRepo   repos.ImageRepository
	programRepo repos.ProgramRepository
	recRepo     repos.RecordingRepository
}

type PandaProgramExecutorJob struct {
	imageRepo  repos.ImageRepository
	recRepo    repos.RecordingRepository
	opts       *PandaProgramExecutorOptions
	agent      *dockerPandaAgent
	Recordings []PandaAgentRecording
}

type PandaProgramExecutorOptions struct {
	Name    *string
	Image   *models.Image
	Program *models.InteractionProgram
}

func NewPandaProgramExecutor(ctx context.Context) (*PandaProgramExecutor, error) {
	imageRepo, err := repos.GetImageRepository(ctx)
	if err != nil {
		return nil, err
	}

	programRepo, err := repos.GetProgramRepository(ctx)
	if err != nil {
		return nil, err
	}

	recRepo, err := repos.GetRecordingRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &PandaProgramExecutor{
		imageRepo:   imageRepo,
		programRepo: programRepo,
		recRepo:     recRepo,
	}, nil
}

func (p *PandaProgramExecutor) NewExecutorJob(ctx context.Context, opts *PandaProgramExecutorOptions) (*PandaProgramExecutorJob, error) {
	// Basic rundown of what will happen:
	// 1. pull information from the database
	//    - making this the caller's responsibility
	//	  - passed in as 'opts'
	// 2. open a stream to the file in blob storage
	//    - caller's responsibility for now
	// Rest in StartJob

	job := &PandaProgramExecutorJob{
		imageRepo:  p.imageRepo,
		recRepo:    p.recRepo,
		opts:       opts,
		Recordings: []PandaAgentRecording{},
	}
	return job, nil
}

func (p *PandaProgramExecutorJob) setupContainer(ctx context.Context) error {
	// 3. create a panda instance using that file
	agent, err := CreateDockerPandaAgent2(ctx)
	if err != nil {
		// todo: return via a channel
		return errors.Wrap(err, "failed to create panda agent")
	}
	p.agent = agent

	err = p.agent.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to connect to agent")
	}

	// Copy the image to the agent
	err = p.copyImageFiles(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to copy image files to agent")
	}

	// 4. start the agent with the given image and configuration
	fmt.Println("Starting PANDA instance")
	err = agent.StartAgentWithOpts(ctx, &StartAgentRequest{
		Config: &p.opts.Image.Config.PandaConfig,
	})
	if err != nil {
		return errors.Wrap(err, "failed to start agent")
	}

	return err
}

func (p *PandaProgramExecutorJob) copyImageFiles(ctx context.Context) error {
	image := p.opts.Image
	for _, file := range p.opts.Image.Files {
		fmt.Printf("Copying file %s to agent\n", file.FileName)
		fileReader, err := p.imageRepo.OpenImageFile(ctx, image.ID, file.ID)
		if err != nil {
			return errors.Wrap(err, "failed to open image file")
		}

		err = p.agent.CopyFileToContainer(ctx, fileReader, file.Size, file.FileName)
		if err != nil {
			return errors.Wrap(err, "failed to copy file to agent")
		}
	}

	return nil
}

func (p *PandaProgramExecutorJob) Run(ctx context.Context) {
	// Initialize the container
	err := p.setupContainer(ctx)
	if err != nil {
		panic(err)
	}

	// Parse the interaction program instructions and run them
	err = p.runProgram(ctx, p.opts.Program)
	if err != nil {
		panic(err)
	}

	// Upload recordings
	err = p.uploadRecordings(ctx)
	if err != nil {
		panic(err)
	}

	// stop the agent
	err = p.agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

}

func (p *PandaProgramExecutorJob) runProgram(ctx context.Context, prog *models.InteractionProgram) error {
	instructions, err := models.ParseInteractionProgram(p.opts.Program.Instructions)
	if err != nil {
		return err
	}

	for _, inst := range instructions {
		err := p.runCommand(ctx, inst)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Done running commands\n")

	return nil
}

func (p *PandaProgramExecutorJob) runCommand(ctx context.Context, cmd models.InteractionProgramInstruction) error {
	if cmd != nil {
		// todo: I think we should make this polymorphic
		fmt.Printf("Running command: %+v", cmd)
		switch strings.ToLower(cmd.GetInstructionType()) {
		case "start_recording":
			// Since we have a start recording command, we have to type cast cmd to a pointer for a StartRecordingInstruction from the models package
			err := p.agent.StartRecording(ctx, cmd.(*models.StartRecordingInstruction).RecordingName)
			if err != nil {
				panic(err)
			}
		case "stop_recording":
			recording, err := p.agent.StopRecording(ctx)
			// push the recording to the recordings array
			p.Recordings = append(p.Recordings, recording)
			if err != nil {
				panic(err)
			}
		case "cmd":
			cmdResult, err := p.agent.RunCommand(ctx, cmd.(*models.RunCommandInstruction).Command)
			if err != nil {
				panic(err)
			}
			fmt.Printf(" %s\n", cmdResult.Logs)
			// result = append(result, cmdResult.Logs+"\n")
		case "filesystem":
			// for the future
		case "network":
			// for the future
		default:
			fmt.Printf("Incorrect Command Type %s, Correct options can be found in the commands.md file\n", cmd.GetInstructionType())
		}
	}
	return nil
}

func (p *PandaProgramExecutorJob) uploadRecordings(ctx context.Context) error {
	for _, rec := range p.Recordings {
		fmt.Println("Uploading recording: " + rec.Name())
		newRecording := &models.Recording{
			ID:          nil,
			ImageID:     p.opts.Image.ID,
			ProgramID:   p.opts.Program.ID,
			Name:        *p.opts.Name + "-" + rec.Name(),
			Description: "",
			Date:        time.Now().String(),
		}

		newRecording, err := p.recRepo.CreateRecording(ctx, newRecording)
		if err != nil {
			return errors.Wrap(err, "failed to create recording")
		}

		ndlogRecordingFile, err := p.recRepo.CreateRecordingFile(ctx, &models.CreateRecordingFileRequest{
			RecordingID: newRecording.ID,
			Name:        *p.opts.Name + "-" + rec.NdlogFilename(),
			FileType:    "ndlog",
		})
		if err != nil {
			return errors.Wrap(err, "failed to create ndlog object")
		}

		// Upload the ndlog file
		ndlogStream, err := rec.OpenNdlog(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to open ndlog file")
		}
		defer ndlogStream.Close()

		_, err = p.recRepo.UploadRecordingFile(ctx, &models.UploadRecordingFileRequest{
			RecordingID: newRecording.ID,
			FileID:      ndlogRecordingFile.ID,
		}, ndlogStream)
		if err != nil {
			return errors.Wrap(err, "failed to upload ndlog object")
		}

		snapshotRecordingFile, err := p.recRepo.CreateRecordingFile(ctx, &models.CreateRecordingFileRequest{
			RecordingID: newRecording.ID,
			Name:        *p.opts.Name + "-" + rec.SnapshotFilename(),
			FileType:    "snapshot",
		})
		if err != nil {
			return errors.Wrap(err, "failed to create snapshot object")
		}

		// Upload the snapshot file
		snapshotStream, err := rec.OpenSnapshot(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to open snapshot file")
		}
		defer snapshotStream.Close()

		_, err = p.recRepo.UploadRecordingFile(ctx, &models.UploadRecordingFileRequest{
			RecordingID: newRecording.ID,
			FileID:      snapshotRecordingFile.ID,
		}, snapshotStream)
		if err != nil {
			return errors.Wrap(err, "failed to upload snapshot object")
		}
	}
	return nil
}
