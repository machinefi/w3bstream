package wasmapi

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ParamValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("eventType") == "" {
			c.Abort()
			c.Error(errors.New("param illegal"))
			return
		}
		c.Next()
	}
}
