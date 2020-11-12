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

	srv := ms.NewMicroService(
		ms.WithGin(handler).WithPrometheus().Listen(5000),
	)

	srv.Start()

}

func handler(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
