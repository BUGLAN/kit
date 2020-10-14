package ms

import (
	"context"
	"fmt"
	_ "github.com/BUGLAN/kit/logutil"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gin-gonic/gin"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/mkevac/debugcharts"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"time"
)

type MicroService struct {
	ctx                    context.Context
	grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamInterceptors []grpc.StreamServerInterceptor
	grpcPreprocess         []func(srv *grpc.Server)
	enableHTTP             bool
	enableGRPC             bool
	engine                 *gin.Engine
	enableHttpCORS         bool
	startTime              time.Time
	logger                 zerolog.Logger
}

type MicroServiceOption func(ms *MicroService)

func NewMicroService(opts ...MicroServiceOption) *MicroService {
	ms := &MicroService{
		logger: log.With().Str("component", "ms").Caller().Logger(),
	}
	for _, opt := range opts {
		opt(ms)
	}
	return ms
}

func WithGinHttpServer(engine *gin.Engine) MicroServiceOption {
	return func(ms *MicroService) {
		ms.engine = engine
	}
}

func WithPrometheus() MicroServiceOption {
	return func(ms *MicroService) {
		// grpc metrics
		ms.grpcUnaryInterceptors = append(ms.grpcUnaryInterceptors, grpcprometheus.UnaryServerInterceptor)
		ms.grpcStreamInterceptors = append(ms.grpcStreamInterceptors, grpcprometheus.StreamServerInterceptor)
		ms.grpcPreprocess = append(ms.grpcPreprocess, func(srv *grpc.Server) {
			grpcprometheus.Register(srv)
		})

		// http metrics
		ms.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}
}

func WithGRPC(preprocess func(srv *grpc.Server)) MicroServiceOption {
	return func(ms *MicroService) {
		ms.grpcPreprocess = append(ms.grpcPreprocess, preprocess)
	}
}

func (ms *MicroService) ListenAndServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		ms.logger.Panic().Err(err).Msg("create net listener fail")
	}

	grpcServer := grpc.NewServer()
	for _, f := range ms.grpcPreprocess {
		f(grpcServer)
	}

	ms.grpcServer(grpcServer, listener)
	reflection.Register(grpcServer)
	ms.httpServer(listener)
	ms.forever()
}

func (ms *MicroService) httpServer(listener net.Listener) {
	go func() {
		if err := ms.engine.RunListener(listener); err != nil {
			ms.logger.Panic().Err(err).Msg("run listener fail")
		}
	}()
}

func (ms *MicroService) grpcServer(grpcServer *grpc.Server, listener net.Listener) {
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			ms.logger.Panic().Err(err).Msg("grpc server fail")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer func() {
			cancel()
		}()

		conn, err := grpc.DialContext(ctx, "127.0.0.1:5000")
		if err != nil {
			ms.logger.Panic().Err(err).Msg("dial grpc server fail")
		}

		defer func() {
			conn.Close()
		}()

		grpcUIHandler, err := standalone.HandlerViaReflection(ms.ctx, conn, listener.Addr().String())
		if err != nil {
			ms.logger.Panic().Err(err).Msg("enable grpcui fail")
		}

		ms.engine.GET("/debug/grpc-ui", gin.WrapH(grpcUIHandler))
	}()
}

func (ms *MicroService) forever() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ms.shutdown()
}

func (ms *MicroService) shutdown() {
	ms.logger.Info().Msg("ms shutdown")
}
