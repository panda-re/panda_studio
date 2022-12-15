package panda_controller

import (
	"context"

	pb "github.com/panda-re/panda_studio/panda_agent/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcPandaAgent struct {
	cc         *grpc.ClientConn
	cli pb.PandaAgentClient
}

const DEFAULT_GRPC_ADDR = "localhost:50051"

func CreateDefaultGrpcPandaAgent() (PandaAgent, error) {
	conn, err := grpc.Dial(DEFAULT_GRPC_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewPandaAgentClient(conn)

	return &grpcPandaAgent{
		cc:         conn,
		cli: client,
	}, nil
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


func (pa *grpcPandaAgent) Close() error {
	return pa.cc.Close()
}
