package main

import (
	"github.com/BUGLAN/kit/ms"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// config parser, (struct, env, yaml, toml)
	// init middleware, (support mysql, redis, mongo, es)
	// config micro service middleware, like trace, monitor， metrics(jaeger, prometheus, grafana, grpc)
	// start micro service

	engine := buildEngine()
	server := ms.NewMicroService(
		ms.WithGinHttpServer(engine),
		ms.WithPrometheus(),
		)
	server.ListenAndServer(5000)
}

func buildEngine() *gin.Engine {
	engine := gin.Default()
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return engine
}
