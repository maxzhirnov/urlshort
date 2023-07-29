package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/mocks"
	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mockURLShortener = mocks.MockURLShortenerService{
	CreateFunc: func(originalURL string) (id string, err error) {
		return "", nil
	},
	GetFunc: func(id string) (url models.URL, err error) {
		return models.URL{}, nil
	},
}

func Test_handleCreate(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name       string
		method     string
		url        string
		createFunc func(string) (string, error)
		want       want
	}{
		{
			name:       "test success",
			method:     http.MethodPost,
			url:        "https://newsite.com",
			createFunc: func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				body:        "http://example.com/12345678",
			},
		},
		{
			name:       "test method GET",
			method:     http.MethodGet,
			url:        "https://newsite.com",
			createFunc: func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "only POST requests allowed",
			},
		},
		{
			name:       "test method GET",
			method:     http.MethodGet,
			url:        "https://newsite.com",
			createFunc: func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "only POST requests allowed",
			},
		},
		{
			name:       "test method DELETE",
			method:     http.MethodDelete,
			url:        "https://newsite.com",
			createFunc: func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "only POST requests allowed",
			},
		},
		{
			name:       "test error",
			method:     http.MethodPost,
			url:        "https://newsite.com",
			createFunc: func(s string) (string, error) { return "", errors.New("error occurred") },
			want: want{
				statusCode:  http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
				body:        "error creating shorten url",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(tt.method, "/", strings.NewReader(tt.url))
			m := &mockURLShortener
			m.CreateFunc = tt.createFunc
			h := handleCreate(m)
			h(c)

			res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.body, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_handleRedirectToOriginal(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		method  string
		reqURL  string
		getFunc func(id string) (models.URL, error)
		want    want
	}{
		{
			name:    "success test case",
			method:  http.MethodGet,
			reqURL:  "/12345678",
			getFunc: func(id string) (models.URL, error) { return models.URL{OriginalURL: "ya.ru", ID: "12345678"}, nil },
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://ya.ru",
			},
		},
		{
			name:    "test wrong method - POST",
			method:  http.MethodPost,
			reqURL:  "/12345678",
			getFunc: func(id string) (models.URL, error) { return models.URL{OriginalURL: "ya.ru", ID: "12345678"}, nil },
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
			},
		},
		{
			name:    "test case error",
			method:  http.MethodGet,
			reqURL:  "/12345678",
			getFunc: func(id string) (models.URL, error) { return models.URL{}, errors.New("error occurred") },
			want: want{
				statusCode: http.StatusNotFound,
				location:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(tt.method, tt.reqURL, nil)
			m := &mockURLShortener
			m.GetFunc = tt.getFunc
			h := handleRedirectToOriginal(m)
			h(c)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}
