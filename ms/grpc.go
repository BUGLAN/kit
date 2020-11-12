package ms

import (
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

const defaultTimeout = time.Second * 5

type GRPC struct {
	enable                 bool
	port                   int
	enablePrometheus       bool
	grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamInterceptors []grpc.StreamServerInterceptor
	grpcPreprocess         []func(srv *grpc.Server)
	grpcUiPort             int
	ui                     *GRPCUI
}

type GRPCUI struct {
	enable bool
	port   int
	mux    *http.ServeMux
}

func WithGRPC(preprocess ...func(srv *grpc.Server)) *GRPC {
	g := &GRPC{
		enable:         true,
		grpcPreprocess: make([]func(srv *grpc.Server), 0),
	}
	g.grpcPreprocess = append(g.grpcPreprocess, preprocess...)
	return g
}

func (g *GRPC) Listen(port int) MicroServiceOption {
	return func(ms *MicroService) {
		g.port = port
		ms.grpc = g
	}
}

func (g *GRPC) WithPrometheus() *GRPC {
	g.enablePrometheus = true
	g.grpcUnaryInterceptors = append(g.grpcUnaryInterceptors, grpcprometheus.UnaryServerInterceptor)
	g.grpcStreamInterceptors = append(g.grpcStreamInterceptors, grpcprometheus.StreamServerInterceptor)
	g.grpcPreprocess = append(g.grpcPreprocess, func(srv *grpc.Server) {
		grpcprometheus.Register(srv)
	})
	return g
}

func (g *GRPC) WithGRPCUI(port int) *GRPC {
	g.ui = &GRPCUI{
		enable: true,
		port:   port,
		mux:    http.NewServeMux(),
	}
	return g
}
