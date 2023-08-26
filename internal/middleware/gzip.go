package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var gzWriter *gzip.Writer

func init() {
	var err error
	gzWriter, err = gzip.NewWriterLevel(nil, gzip.BestSpeed)
	if err != nil {
		log.Fatal(err)
	}
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

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decoding
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

		// Encoding
		contentType := c.GetHeader("Content-Type")
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			if strings.Contains(contentType, gin.MIMEJSON) || strings.Contains(contentType, gin.MIMEHTML) {
				c.Writer = &gzipResponseWriter{
					writer:         gzWriter,
					ResponseWriter: c.Writer,
				}
				c.Writer.Header().Set("Content-Encoding", "gzip")
			}
		}
		c.Next()
	}
}
