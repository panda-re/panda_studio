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

var _ PandaAgent = &grpcPandaAgent{}

const DEFAULT_GRPC_ADDR = "localhost:50051"

// PandaAgent interface
func CreateGrpcPandaAgent(endpoint string) (interface{}, error) {
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

// PandaAgent interface
func CreateDefaultGrpcPandaAgent() (interface{}, error) {
	return CreateGrpcPandaAgent(DEFAULT_GRPC_ADDR)
}

// StartAgent implements PandaAgent
func (pa *grpcPandaAgent) StartAgent(ctx context.Context) error {
	return pa.StartAgentWithOpts(ctx, &pb.StartAgentRequest{Config: &DEFAULT_CONFIG})
}

// StartAgentWithOpts implements PandaAgent
// Uses x86_64 generic defaults
func (pa *grpcPandaAgent) StartAgentWithOpts(ctx context.Context, opts *pb.StartAgentRequest) error {
	_, err := pa.cli.StartAgent(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}

// StopAgent implements PandaAgent
func (pa *grpcPandaAgent) StopAgent(ctx context.Context) error {
	_, err := pa.cli.StopAgent(ctx, &pb.StopAgentRequest{})
	if err != nil {
		return err
	}

	return nil
}

// RunCommand implements PandaAgent
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
func (pa *grpcPandaAgent) StopRecording(ctx context.Context) (PandaAgentRecording, error) {
	resp, err := pa.cli.StopRecording(ctx, &pb.StopRecordingRequest{})
	if err != nil {
		return nil, err
	}

	return &GenericPandaAgentRecordingConcrete{
		RecordingName: resp.RecordingName,
	}, nil
}

// StartReplay implements PandaAgent
// Uses x86_64 generic defaults
func (pa *grpcPandaAgent) StartReplay(ctx context.Context, recordingName string) (*PandaAgentReplayResult, error) {
	return pa.StartReplayWithOpts(ctx, &pb.StartAgentRequest{Config: &DEFAULT_CONFIG}, recordingName)
}

// StartReplayWithOpts implements PandaAgent
func (pa *grpcPandaAgent) StartReplayWithOpts(ctx context.Context, opts *pb.StartAgentRequest, recordingName string) (*PandaAgentReplayResult, error) {
	resp, err := pa.cli.StartReplay(ctx, &pb.StartReplayRequest{
		Config:        opts.Config,
		RecordingName: recordingName,
	})
	if err != nil {
		return nil, err
	}

	return &PandaAgentReplayResult{
		Serial: resp.GetSerial(),
		Replay: resp.GetReplay(),
	}, nil
}

// StopReplay implements PandaAgent
func (pa *grpcPandaAgent) StopReplay(ctx context.Context) (*PandaAgentReplayResult, error) {
	resp, err := pa.cli.StopReplay(ctx, &pb.StopReplayRequest{})
	if err != nil {
		return nil, err
	}
	return &PandaAgentReplayResult{
		Serial: resp.GetSerial(),
		Replay: resp.GetReplay(),
	}, nil
}

func (pa *grpcPandaAgent) Close() error {
	return pa.cc.Close()
}

func (pa *grpcPandaAgent) SendNetworkCommand(ctx context.Context, in *NetworkRequest, opts ...grpc.CallOption) (*NetworkResponse, error) {
	resp, err := pa.cli.SendNetworkCommand(ctx, &pb.NetworkRequest{
		SocketType:   in.SocketType,
		Port:         in.Port,
		Application:  in.Application,
		Command:      in.Command,
		CustomPacket: in.CustomPacket,
	})

	if err != nil {
		return nil, err
	}

	return &NetworkResponse{
		StatusCode: resp.GetStatusCode(),
		Output:     resp.GetOutput(),
	}, nil
}
