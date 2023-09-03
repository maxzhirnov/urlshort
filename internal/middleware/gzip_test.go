package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (l mockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l mockLogger) Error(msg string, keysAndValues ...interface{}) {}
func (l mockLogger) Fatal(msg string, keysAndValues ...interface{}) {}

func TestGzipMiddleware(t *testing.T) {
	gzipWriter, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
	if err != nil {
		t.Error(err)
	}
	r := gin.Default()
	r.Use(Gzip(&mockLogger{}, gzipWriter))

	r.POST("/test", func(c *gin.Context) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		assert.Equal(t, "test body", buf.String())

		c.String(http.StatusOK, "response body")
	})

	// Insert a gzip compressed request body
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte("test body")); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	// Send gzipped request
	req, _ := http.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response is gzipped
	gr, err := gzip.NewReader(w.Body)
	assert.Nil(t, err)
	respBody, err := io.ReadAll(gr)
	respBodyString := string(respBody)
	assert.Nil(t, err)
	assert.Equal(t, "response body", respBodyString)
}
