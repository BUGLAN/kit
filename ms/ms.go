package ms

import (
	"context"
	"fmt"
	"github.com/BUGLAN/kit/config"
	"github.com/BUGLAN/kit/logutil"
	_ "github.com/BUGLAN/kit/logutil"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gin-gonic/gin"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/mkevac/debugcharts"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const defaultTimeout = 5

type MicroService struct {
	ctx                    context.Context
	grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamInterceptors []grpc.StreamServerInterceptor
	grpcPreprocess         []func(srv *grpc.Server)
	enableHTTP             bool
	enableGRPC             bool
	debug                  bool
	engine                 *gin.Engine
	enableHttpCORS         bool
	startTime              time.Time
	logger                 zerolog.Logger
	grpcPort               int
	httpPort               int
	mux                    *http.ServeMux
	config                 *config.KitConfig
}

type MicroServiceOption func(ms *MicroService)

func NewMicroService(opts ...MicroServiceOption) *MicroService {
	ms := &MicroService{
		logger: logutil.NewLogger("ms.component", "ms"),
		ctx:    context.Background(),
		mux:    http.NewServeMux(),
		config: config.NewKitConfig(),
	}

	if *config.Debug {
		ms.debug = true
		ms.logger.Info().Msg("Enable debug mode")
		logutil.SetGlobalLevel("debug")
	}

	for _, opt := range opts {
		opt(ms)
	}
	return ms
}

func WithGinHTTP(handler func(engine *gin.Engine)) MicroServiceOption {
	return func(ms *MicroService) {
		ms.enableHTTP = true
		ms.engine = gin.Default()
		handler(ms.engine)
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

func WithGRPC(port int, preprocess func(srv *grpc.Server)) MicroServiceOption {
	return func(ms *MicroService) {
		ms.enableGRPC = true
		ms.grpcPort = port
		ms.grpcPreprocess = append(ms.grpcPreprocess, preprocess)
	}
}

func (ms *MicroService) ListenAndServer(port int) {
	if ms.enableGRPC {
		ms.grpcServer()
	}

	// http server
	go func() {
		ms.engine.Run(fmt.Sprintf(":%d", port))
	}()

	go func() {
		if ms.enableGRPC && ms.debug {
			ms.enableDebug()
			http.ListenAndServe(":8888", ms.mux)
		}
	}()

	ms.forever()
}

// func (ms *MicroService) httpServer(listener net.Listener) {
//	go func() {
//		if err := ms.engine.RunListener(listener); err != nil {
//			ms.logger.Panic().Err(err).Msg("run listener fail")
//		}
//	}()
// }

func (ms *MicroService) grpcServer() {
	s := grpc.NewServer()
	for _, f := range ms.grpcPreprocess {
		f(s)
	}
	reflection.Register(s)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ms.grpcPort))
	if err != nil {
		ms.logger.Panic().Err(err).Msg("create net listener fail")
	}
	go func() {
		err := s.Serve(listener)
		if err != nil {
			ms.logger.Panic().Err(err).Msg("grpc server fail")
		}
	}()
}

func (ms *MicroService) enableDebug() {
	addr := fmt.Sprintf(":%d", ms.grpcPort)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*defaultTimeout)
		defer func() {
			cancel()
		}()

		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
		if err != nil {
			ms.logger.Panic().Err(err).Msg("dial grpc server fail")
		}
		defer func() {
			conn.Close()
		}()

		h, err := standalone.HandlerViaReflection(ms.ctx, conn, addr)
		if err != nil {
			ms.logger.Panic().Err(err).Msg("enable grpcui fail")
		}
		ms.logger.Info().Msg("enable grpc ui")
		ms.mux.Handle("/grpcui/", http.StripPrefix("/grpcui", h))
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
