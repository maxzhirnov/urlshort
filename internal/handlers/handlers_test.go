package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockURLShortenerService struct {
	CreateFunc func(originalURL string) (id string, err error)
	GetFunc    func(id string) (url models.URL, err error)
}

func (m *mockURLShortenerService) Create(originalURL string) (string, error) {
	return m.CreateFunc(originalURL)
}

func (m *mockURLShortenerService) Get(id string) (models.URL, error) {
	return m.GetFunc(id)
}

var mockURLShortener = mockURLShortenerService{
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
		name         string
		method       string
		url          string
		redirectHost string
		createFunc   func(string) (string, error)
		want         want
	}{
		{
			name:         "test success without protocol",
			method:       http.MethodPost,
			url:          "newsite.com",
			redirectHost: "https://example.com",
			createFunc:   func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				body:        "https://example.com/12345678",
			},
		},
		{
			name:         "test success with protocol http",
			method:       http.MethodPost,
			url:          "http://newsite.com",
			redirectHost: "https://example.com",
			createFunc:   func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				body:        "https://example.com/12345678",
			},
		},
		{
			name:         "test localhost:8080",
			method:       http.MethodPost,
			url:          "http://newsite.com",
			redirectHost: "localhost:8080",
			createFunc:   func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				body:        "localhost:8080/12345678",
			},
		},
		{
			name:         "test success with protocol https",
			method:       http.MethodPost,
			url:          "http://newsite.com",
			redirectHost: "https://example.com",
			createFunc:   func(s string) (string, error) { return "12345678", nil },
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				body:        "https://example.com/12345678",
			},
		},
		{
			name:         "test error",
			method:       http.MethodPost,
			url:          "https://newsite.com",
			redirectHost: "https://example.com",
			createFunc:   func(s string) (string, error) { return "", errors.New("error occurred") },
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
			handlers := NewShortenerHandlers(m, tt.redirectHost)
			h := handlers.HandleCreate()
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

func Test_handleRedirect(t *testing.T) {
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
				location:   "http://ya.ru",
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
			handlers := NewShortenerHandlers(m, "")
			h := handlers.HandleRedirect()
			h(c)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}