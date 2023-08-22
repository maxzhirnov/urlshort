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

func TestGzipMiddleware(t *testing.T) {
	r := gin.Default()
	r.Use(GzipMiddleware())

	r.POST("/test", func(c *gin.Context) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		assert.Equal(t, "test body", buf.String())

		c.String(http.StatusOK, "response body")
	})

	// Create a gzip compressed request body
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
	assert.Nil(t, err)
	assert.Equal(t, "response body", string(respBody))
}
