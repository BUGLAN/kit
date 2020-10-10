package main

import "github.com/BUGLAN/kit/ms"

func main() {
	// config parser, (struct, env, yaml, toml)
	// init middleware, (support mysql, redis, mongo, es)
	// config micro service middleware, like trace, monitor， metrics(jaeger, prometheus, grafana, grpc)
	// start micro service

	server := ms.NewMicroService()
	server.ListenAndServer(5000)
}