package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Hello(ctx *gin.Context) {
	name := ctx.Request.Header.Get("name")
	fmt.Fprintf(ctx.Writer, "hello %v", name)
}
