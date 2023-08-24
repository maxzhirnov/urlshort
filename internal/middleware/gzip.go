package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

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

		// Упаковка
		contentType := c.GetHeader("Content-Type")
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			if strings.Contains(contentType, gin.MIMEJSON) || strings.Contains(contentType, gin.MIMEHTML) {
				gzWriter := gzip.NewWriter(c.Writer)
				defer gzWriter.Close()

				c.Writer.Header().Set("Content-Encoding", "gzip")
				c.Writer = &gzipResponseWriter{
					writer:         gzWriter,
					ResponseWriter: c.Writer,
				}
			}
		}
		c.Next()
	}
}
