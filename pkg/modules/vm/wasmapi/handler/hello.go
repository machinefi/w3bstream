package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Hello(c *gin.Context) {
	name := c.Request.Header.Get("name")

	if err := h.setAsync(c); err != nil {
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf(`{"message": "hello %s"}`, name))
}

func (h *Handler) HelloAsync(c *gin.Context) {
	name := c.Request.Header.Get("name")
	c.JSON(http.StatusOK, fmt.Sprintf(`{"async message": "hello %s"}`, name))
}
