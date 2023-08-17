package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/logging"
	"time"
)

func Logging(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		resStatusCode := c.Writer.Status()
		resContentLength := c.Writer.Size()

		logger.Info("API request",
			"method", reqMethod,
			"uri", reqURI,
			"latency", latencyTime.Seconds(),
		)

		logger.Info("API response",
			"status_code", resStatusCode,
			"content_length", resContentLength)
	}
}
