package main

import (
	"context"
	"github.com/BUGLAN/kit/internal/example/grpc-example/pb"
	"github.com/BUGLAN/kit/ms"
	"google.golang.org/grpc"
)

type GreetService struct {
}

func NewGreetService() *GreetService {
	return &GreetService{}
}

func (g GreetService) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Msg: "PONG"}, nil
}

func main() {
	s := NewGreetService()
	srv := ms.NewMicroService(
		ms.WithGRPC(func(srv *grpc.Server) {
			pb.RegisterGreetServer(srv, s)
		}).WithPrometheus().WithGRPCUI(5001).Listen(5000),
	)
	srv.Start()
}
