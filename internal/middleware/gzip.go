package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
}

type gzipResponseWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	g.writer.Reset(g.ResponseWriter)
	defer g.writer.Close()
	return g.writer.Write(data)
}

func Gzip(logger logger, gzipWriter *gzip.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decoding
		if c.GetHeader("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				logger.Error(err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer gr.Close()

			bodyBytes, err := io.ReadAll(gr)
			if err != nil {
				logger.Error(err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Encoding
		contentType := c.GetHeader("Content-Type")
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			if strings.Contains(contentType, gin.MIMEJSON) || strings.Contains(contentType, gin.MIMEHTML) {
				logger.Info("gziping data")
				c.Writer = &gzipResponseWriter{
					writer:         gzipWriter,
					ResponseWriter: c.Writer,
				}
				c.Writer.Header().Set("Content-Encoding", "gzip")
			}
		}
		c.Next()
	}
}
