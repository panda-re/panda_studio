package panda_controller

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/db/repos"
	"github.com/pkg/errors"
)

type PandaProgramExecutor struct {
	imageRepo repos.ImageRepository
	programRepo repos.ProgramRepository
}

type PandaProgramExecutorJob struct {
	imageRepo repos.ImageRepository
	opts *PandaProgramExecutorOptions
	agent *DockerGrpcPandaAgent2
	recordings []PandaAgentRecording
}

type PandaProgramExecutorOptions struct {
	Image *models.Image
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

	return &PandaProgramExecutor{
		imageRepo: imageRepo,
		programRepo: programRepo,
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
		imageRepo: p.imageRepo,
		opts: opts,
		recordings: []PandaAgentRecording{},
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

	// stop the agent
	err = p.agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}

	// 7. upload the recording files to blob storage
	for _, recording := range p.recordings {
		fmt.Printf("Copying recording %s to local disk\n", recording.Name())
		ndlogStream, err := recording.OpenNdlog(ctx)
		if err != nil {
			panic(err)
		}
		defer ndlogStream.Close()
		ndlogFile, err := os.OpenFile(recording.NdlogFilename(), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer ndlogFile.Close()
		_, err = io.Copy(ndlogFile, ndlogStream)
		if err != nil {
			panic(err)
		}
		snpStream, err := recording.OpenNdlog(ctx)
		if err != nil {
			panic(err)
		}
		defer snpStream.Close()
		snpFile, err := os.OpenFile(recording.SnapshotFilename(), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer snpFile.Close()
		_, err = io.Copy(snpFile, snpStream)
		if err != nil {
			panic(err)
		}
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

	return nil
}

func (p *PandaProgramExecutorJob) runCommand(ctx context.Context, cmd models.InteractionProgramInstruction) error {
	if cmd != nil {
		// todo: I think we should make this polymorphic
		switch cmd.GetInstructionType() {
		case "start_recording":
			// Since we have a start recording command, we have to type cast cmd to a pointer for a StartRecordingInstruction from the models package
			err := p.agent.StartRecording(ctx, cmd.(*models.StartRecordingInstruction).RecordingName)
			if err != nil {
				panic(err)
			}
			break
		case "stop_recording":
			recording, err := p.agent.StopRecording(ctx)
			// push the recording to the recordings array
			p.recordings = append(p.recordings, recording)
			if err != nil {
				panic(err)
			}
			break
		case "command":
			cmdResult, err := p.agent.RunCommand(ctx, cmd.(*models.RunCommandInstruction).Command)
			if err != nil {
				panic(err)
			}
			fmt.Printf(" %s\n", cmdResult.Logs)
			// result = append(result, cmdResult.Logs+"\n")
			break
		case "filesystem":
			// for the future
		case "network":
			// for the future
		default:
			fmt.Printf("Incorrect Command Type, Correct options can be found in the commands.md file")
			break
		}
	}
	return nil
}
