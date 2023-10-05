package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger logger) gin.HandlerFunc {
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
