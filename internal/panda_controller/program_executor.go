package panda_controller

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/panda-re/panda_studio/internal/db/models"
	"github.com/panda-re/panda_studio/internal/db/repos"
	"github.com/panda-re/panda_studio/panda_agent/pb"
)

type PandaProgramExecutor struct {
	imageRepo repos.ImageRepository
}

type PandaProgramExecutorJob struct {
	opts *PandaProgramExecutorOptions
	agent *DockerGrpcPandaAgent2
	recordings []*PandaAgentRecording
}

type PandaProgramExecutorOptions struct {
	Image *models.Image
	Program *models.InteractionProgram
	Instructions models.InteractionProgramInstructionList
	ImageFileReader io.ReadSeeker
}

func (p *PandaProgramExecutor) NewExecutorJob(opts *PandaProgramExecutorOptions) (*PandaProgramExecutorJob, error) {
	// Basic rundown of what will happen:
	// 1. pull information from the database
	//    - making this the caller's responsibility
	//	  - passed in as 'opts'
	// 2. open a stream to the file in blob storage
	//    - caller's responsibility for now
	// Rest in StartJob

	job := &PandaProgramExecutorJob{
		opts: opts,
		recordings: []*PandaAgentRecording{},
	}
	return job, nil
}

func getReaderSize(reader io.ReadSeeker) (size int64, err error) {
	// Get the size of the file
	// https://stackoverflow.com/a/24563853
	currentPos, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	defer reader.Seek(currentPos, io.SeekStart)

	endPos, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	defer reader.Seek(currentPos, io.SeekStart)

	return endPos, nil
}

func (p *PandaProgramExecutorJob) StartJob(ctx context.Context) {
	// 3. create a panda instance using that file
	agent, err := CreateDockerPandaAgent2(ctx)
	if err != nil {
		// todo: return via a channel
		panic(err)
	}
	p.agent = agent

	err = p.agent.Connect(ctx)
	if err != nil {
		panic(err)
	}

	// Copy the image to the agent
	fmt.Println("Copying image to agent...")
	fileSize, err := getReaderSize(p.opts.ImageFileReader)
	if err != nil {
		panic(err)
	}
	err = p.agent.CopyFileToContainer(ctx, p.opts.ImageFileReader, fileSize, "system_image.qcow2")
	if err != nil {
		panic(err)
	}

	// 4. start the agent with the given image and configuration
	fmt.Println("Starting PANDA instance")
	err = agent.StartAgentWithOpts(ctx, &StartAgentRequest{
		Config: &pb.PandaConfig{
			QcowFileName: "system_image.qcow2",
			Arch: "x86_64",
			Os: "linux-64-ubuntu:4.15.0-72-generic-noaslr-nokaslr",
			Prompt: "root@ubuntu:.*#",
			Cdrom: "ide1-cd0",
			Snapshot: "root",
			Memory: "1024M",
			ExtraArgs: "-display none",
		},
	})
	if err != nil {
		panic(err)
	}
	// 5. send the commands to the agent
	//    - offer an interface for real-time feedback, even if we don't currently use it
	//    - keep track of any recording files that are created
	for _, inst := range p.opts.Instructions {
		err = p.runCommand(ctx, inst)
		if err != nil {
			panic(err)
		}
	}
	// 6. stop the agent
	err = p.agent.StopAgent(ctx)
	if err != nil {
		panic(err)
	}
	// 7. upload the recording files to blob storage
	for _, recording := range p.recordings {
		fmt.Printf("Copying recording %s to local disk\n", recording.RecordingName)
		ndlogStream, err := p.agent.CopyFileFromContainer(ctx, recording.GetNdlogFileName())
		if err != nil {
			panic(err)
		}
		defer ndlogStream.Close()
		ndlogFile, err := os.OpenFile(recording.GetNdlogFileName(), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer ndlogFile.Close()
		_, err = io.Copy(ndlogFile, ndlogStream)
		if err != nil {
			panic(err)
		}
		snpStream, err := p.agent.CopyFileFromContainer(ctx, recording.GetSnapshotFileName())
		if err != nil {
			panic(err)
		}
		defer snpStream.Close()
		snpFile, err := os.OpenFile(recording.GetSnapshotFileName(), os.O_CREATE|os.O_WRONLY, 0644)
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
