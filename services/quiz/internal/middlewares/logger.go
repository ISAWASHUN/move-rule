package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Info("request",
			"status", status,
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)
	}
}
