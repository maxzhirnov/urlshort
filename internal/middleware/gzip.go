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
		// Only handle when it's application/json or text/html
		if c.GetHeader("Content-Type") != "application/json" && c.GetHeader("Content-Type") != "text/html" {
			c.Next()
			return
		}
		// Decompress if needed
		if c.GetHeader("Content-Encoding") == "gzip" {
			gr, err := NewGzipReader(c.Request.Body)
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

		// Check response compression and set writer
		gw := NewGzipResponseWriter(c.Writer)
		defer gw.Close()

		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Writer.Header().Set("Content-Encoding", "gzip")
			c.Writer = gw
		}

		c.Next()
	}
}

type gzipResponseWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

func NewGzipResponseWriter(w gin.ResponseWriter) *gzipResponseWriter {
	gw := gzip.NewWriter(w)
	return &gzipResponseWriter{gw, w}
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipResponseWriter) Close() error {
	return g.writer.Close()
}

type gzipReader struct {
	*gzip.Reader
}

func NewGzipReader(r io.Reader) (*gzipReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &gzipReader{gr}, nil
}
