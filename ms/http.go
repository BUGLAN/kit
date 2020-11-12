package ms

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTPHandler func(*gin.Engine)

type HTTP struct {
	enable           bool
	port             int
	handler          HTTPHandler
	enablePrometheus bool
	enableHttpCORS   bool
	engine           *gin.Engine
}

func WithGin(handler HTTPHandler) *HTTP {
	return &HTTP{
		handler: handler,
		enable:  true,
		engine:  gin.Default(),
	}
}

// WithPrometheus enable prometheus
func (h *HTTP) WithPrometheus() *HTTP {
	h.enablePrometheus = true
	return h
}

func (h *HTTP) Listen(port int) MicroServiceOption {
	return func(ms *MicroService) {
		h.port = port
		ms.http = h
	}
}

// Done mean finish http server
func (h *HTTP) Done() MicroServiceOption {
	return func(ms *MicroService) {
		ms.http = h
	}
}

// start http server
func (ms *MicroService) httpServer() {
	// 未开启 http 服务
	if ms.http == nil || !ms.http.enable {
		return
	}

	// metrics
	if ms.http.enablePrometheus {
		ms.http.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// handler
	ms.http.handler(ms.http.engine)

	// listen and server
	if err := ms.http.engine.Run(fmt.Sprintf(":%d", ms.http.port)); err != nil {
		ms.logger.Panic().Err(err)
	}
}
