package wasmapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func hello(ctx *gin.Context) {
	name := ctx.Request.Header.Get("name")
	fmt.Fprintf(ctx.Writer, "hello %v", name)
}
