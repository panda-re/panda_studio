package panda_controller

import (
	"context"

	pb "github.com/panda-re/panda_studio/panda_agent/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcPandaAgent struct {
	cc  *grpc.ClientConn
	cli pb.PandaAgentClient
}

const DEFAULT_GRPC_ADDR = "localhost:50051"

func CreateGrpcPandaAgent(endpoint string) (PandaAgent, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewPandaAgentClient(conn)

	return &grpcPandaAgent{
		cc:  conn,
		cli: client,
	}, nil
}

func CreateGrpcPandaReplayAgent(endpoint string) (PandaReplayAgent, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewPandaAgentClient(conn)

	return &grpcPandaAgent{
		cc:  conn,
		cli: client,
	}, nil
}

func CreateDefaultGrpcPandaAgent() (PandaAgent, error) {
	return CreateGrpcPandaAgent(DEFAULT_GRPC_ADDR)
}

func (pa *grpcPandaAgent) StartAgent(ctx context.Context) error {
	_, err := pa.cli.StartAgent(ctx, &pb.StartAgentRequest{})
	if err != nil {
		return err
	}

	return nil
}

func (pa *grpcPandaAgent) StopAgent(ctx context.Context) error {
	_, err := pa.cli.StopAgent(ctx, &pb.StopAgentRequest{})
	if err != nil {
		return err
	}

	return nil
}

func (pa *grpcPandaAgent) RunCommand(ctx context.Context, cmd string) (*PandaAgentRunCommandResult, error) {
	resp, err := pa.cli.RunCommand(ctx, &pb.RunCommandRequest{
		Command: cmd,
	})
	if err != nil {
		return nil, err
	}

	return &PandaAgentRunCommandResult{
		Logs: resp.GetOutput(),
	}, nil
}

// StartRecording implements PandaAgent
func (pa *grpcPandaAgent) StartRecording(ctx context.Context, recordingName string) error {
	_, err := pa.cli.StartRecording(ctx, &pb.StartRecordingRequest{
		RecordingName: recordingName,
	})
	if err != nil {
		return err
	}

	return nil
}

// StopRecording implements PandaAgent
func (pa *grpcPandaAgent) StopRecording(ctx context.Context) (*PandaAgentRecording, error) {
	resp, err := pa.cli.StopRecording(ctx, &pb.StopRecordingRequest{})
	if err != nil {
		return nil, err
	}

	return &PandaAgentRecording{
		RecordingName: resp.RecordingName,
		// We cannot know the location with the information we have
		Location: "??",
	}, nil
}

func (pa *grpcPandaAgent) StartReplay(ctx context.Context, recordingName string) (*PandaAgentRunCommandResult, error) {
	a, err := pa.cli.StartReplay(ctx, &pb.StartReplayRequest{
		RecordingName: recordingName,
	})
	if err != nil {
		return nil, err
	}

	return &PandaAgentRunCommandResult{
		Logs: a.String(),
	}, nil
}

func (pa *grpcPandaAgent) StopReplay(ctx context.Context) error {
	_, err := pa.cli.StopReplay(ctx, &pb.StopReplayRequest{})
	if err != nil {
		return err
	}
	return nil
}

func (pa *grpcPandaAgent) Close() error {
	return pa.cc.Close()
}
