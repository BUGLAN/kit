# kit

> a golang kit tools


## Example

* minimum example

```go
package main

import (
	"github.com/BUGLAN/kit/ms"
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {
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
```


* http restful


```go
package main

func main() {

}
```