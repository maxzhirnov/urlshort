package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Распаковка
		if c.GetHeader("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer gr.Close()

			bodyBytes, err := io.ReadAll(gr)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		gw := gzip.NewWriter(c.Writer)
		defer gw.Close()

		// Прихраним оригинальный writer
		originalWriter := c.Writer
		c.Writer = &gzipResponseWriter{gw, originalWriter}

		c.Next()

		contentType := c.Writer.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
			if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
				c.Writer.Header().Set("Content-Encoding", "gzip")
			}
		} else {
			// Если Content-Type не json или html, восстанавливаем оригинальный writer
			c.Writer = originalWriter
		}
	}
}

type gzipResponseWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}
