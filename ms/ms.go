package ms

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fullstorydev/grpcui/standalone"

	"github.com/BUGLAN/kit/config"
	"github.com/BUGLAN/kit/logutil"
	_ "github.com/BUGLAN/kit/logutil"
	_ "github.com/mkevac/debugcharts"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MicroServiceOption func(ms *MicroService)

type MicroService struct {
	ctx       context.Context
	debug     bool
	startTime time.Time
	logger    zerolog.Logger
	mux       *http.ServeMux
	config    *config.KitConfig
	http      *HTTP
	grpc      *GRPC
}

func NewMicroService(opts ...MicroServiceOption) *MicroService {
	ms := &MicroService{
		logger: logutil.NewLogger("ms.component", "ms"),
		ctx:    context.Background(),
		mux:    http.NewServeMux(),
		config: config.NewKitConfig(),
	}

	if *config.Debug {
		ms.debug = true
		ms.logger.Info().Msg("enable debug mode")
		logutil.SetGlobalLevel("debug")
	}

	for _, opt := range opts {
		opt(ms)
	}
	return ms
}

func (ms *MicroService) grpcServer() {
	if ms.grpc == nil || !ms.grpc.enable {
		return
	}
	s := grpc.NewServer()
	for _, f := range ms.grpc.grpcPreprocess {
		f(s)
	}
	reflection.Register(s)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ms.grpc.port))
	if err != nil {
		ms.logger.Panic().Err(err).Msg("create net listener fail")
		return
	}

	go func() {
		ms.logger.Info().Msgf("Listen grpc server on :%d", ms.grpc.port)
		if err := s.Serve(listener); err != nil {
			ms.logger.Panic().Err(err).Msg("grpc server start fail")
		}
	}()

	go func() {
		if !ms.grpc.ui.enable {
			return
		}
		ms.grpcUiServer()
	}()
}

func (ms *MicroService) grpcUiServer() {
	addr := fmt.Sprintf(":%d", ms.grpc.port)
	ctx, cancel := context.WithTimeout(ms.ctx, time.Second*defaultTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		ms.logger.Panic().Err(err)
	}

	h, err := standalone.HandlerViaReflection(ctx, conn, addr)
	if err != nil {
		ms.logger.Panic().Err(err).Msg("Enable grpcui fail")
	}

	ms.logger.Info().Msg("Enable grpc ui")
	ms.mux.Handle("/grpcui/", http.StripPrefix("/grpcui", h))

	// handler server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", ms.grpc.ui.port), ms.mux); err != nil {
		ms.logger.Panic().Err(err)
	}
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

func (ms *MicroService) Start() {
	// start http server
	go ms.httpServer()

	// start grpc server
	go ms.grpcServer()

	// start worker

	// forever
	ms.forever()
}
