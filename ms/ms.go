package ms

import (
	"context"
	"fmt"
	_ "github.com/BUGLAN/kit/logutil"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"time"
)

type MicroService struct {
	ctx            context.Context
	engine         *gin.Engine
	enableHttpCORS bool
	startTime      time.Time
	logger         zerolog.Logger
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

func (ms *MicroService) ListenAndServer(port int) {

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		ms.logger.Panic().Err(err).Msg("create net listener fail")
	}

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

func (ms *MicroService) forever() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ms.shutdown()
}

func (ms *MicroService) shutdown() {
	ms.logger.Info().Msg("ms shutdown")
}
