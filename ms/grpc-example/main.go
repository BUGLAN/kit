package main

import (
	"context"
	ms2 "github.com/BUGLAN/kit/ms"
	"github.com/BUGLAN/kit/ms/grpc-example/pb"
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
	ms := ms2.NewMicroService(
		ms2.WithGRPC(5000, func(srv *grpc.Server) {
			pb.RegisterGreetServer(srv, s)
		}),
		ms2.WithPrometheus(),
	)
	ms.ListenAndServer(5001)
}