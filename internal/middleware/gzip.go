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

		// Хранение текущего Writer для восстановления позднее
		originalWriter := c.Writer
		// Проверка необходимости сжатия в ответе и установка writer
		gw := NewGzipResponseWriter(c.Writer)
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Writer = gw
		}

		c.Next() // Передача управления следующему обработчику

		// Проверка Content-Type ответа
		contentType := c.Writer.Header().Get("Content-Type")
		if !(strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html")) {
			// Если Content-Type не соответствует ожидаемому, восстанавливаем оригинальный writer и завершаем функцию
			c.Writer = originalWriter
			return
		}

		// Если Content-Type соответствует ожидаемому, устанавливаем заголовок сжатия
		c.Writer.Header().Set("Content-Encoding", "gzip")

		// Закрытие gw
		_ = gw.Close()
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
