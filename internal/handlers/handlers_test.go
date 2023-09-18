package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maxzhirnov/urlshort/internal/auth"
	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/models"
)

type mockURLShortenerService struct {
	CreateFunc func(originalURL string) (url models.ShortURL, err error)
	GetFunc    func(id string) (url models.ShortURL, err error)
}

func (m *mockURLShortenerService) Create(url, uuid string) (models.ShortURL, error) {
	return m.CreateFunc(url)
}

func (m *mockURLShortenerService) CreateBatch(urls []string, uuid string) (ids []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockURLShortenerService) GetAllUsersURLs(uuid string) ([]models.ShortURL, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockURLShortenerService) Get(id string) (models.ShortURL, error) {
	return m.GetFunc(id)
}

func (m *mockURLShortenerService) Ping() error {
	return nil
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
		createFunc   func(string) (models.ShortURL, error)
		want         want
	}{
		{
			name:         "test success without protocol",
			method:       http.MethodPost,
			url:          "newsite.com",
			redirectHost: "https://example.com",
			createFunc: func(s string) (models.ShortURL, error) {
				return models.ShortURL{
					ID:          "12345678",
					OriginalURL: "example.com",
				}, nil
			},
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
			createFunc:   func(s string) (models.ShortURL, error) { return models.ShortURL{ID: "12345678"}, nil },
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
			createFunc:   func(s string) (models.ShortURL, error) { return models.ShortURL{ID: "12345678"}, nil },
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
			createFunc:   func(s string) (models.ShortURL, error) { return models.ShortURL{ID: "12345678"}, nil },
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
			createFunc:   func(s string) (models.ShortURL, error) { return models.ShortURL{}, errors.New("error occurred") },
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
			m := &mockURLShortenerService{}
			m.CreateFunc = tt.createFunc
			lg := logging.NewLogrusLogger(logrus.DebugLevel)
			auths := auth.NewAuth()
			handlers := NewHandlers(m, tt.redirectHost, auths, lg)
			h := handlers.HandleCreate
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
		getFunc func(id string) (models.ShortURL, error)
		want    want
	}{
		{
			name:   "success test case",
			method: http.MethodGet,
			reqURL: "/12345678",
			getFunc: func(id string) (models.ShortURL, error) {
				return models.ShortURL{OriginalURL: "ya.ru", ID: "12345678"}, nil
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "http://ya.ru",
			},
		},
		{
			name:    "test case error",
			method:  http.MethodGet,
			reqURL:  "/12345678",
			getFunc: func(id string) (models.ShortURL, error) { return models.ShortURL{}, errors.New("error occurred") },
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
			m := &mockURLShortenerService{}
			m.GetFunc = tt.getFunc
			handlers := NewHandlers(m, "", nil, nil)
			h := handlers.HandleRedirect
			h(c)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}

func TestHandleShorten(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		input          []byte
		expectedStatus int
		mockFunc       func(originalURL string) (url models.ShortURL, err error)
	}{
		{
			name:           "invalid json",
			input:          []byte(`{"invalid":"json"`),
			expectedStatus: http.StatusBadRequest,
			mockFunc:       mockURLShortenerService{}.CreateFunc,
		},
		{
			name:           "short url",
			input:          []byte(`{"URL": "ab"}`),
			expectedStatus: http.StatusBadRequest,
			mockFunc:       mockURLShortenerService{}.CreateFunc,
		},
		{
			name:           "internal server error",
			input:          []byte(`{"URL": "https://example.com"}`),
			expectedStatus: http.StatusInternalServerError,
			mockFunc: func(originalURL string) (url models.ShortURL, err error) {
				return models.ShortURL{}, errors.New("mocked error")
			},
		},
		{
			name:           "successful shorten",
			input:          []byte(`{"URL": "https://example.com"}`),
			expectedStatus: http.StatusCreated,
			mockFunc: func(originalURL string) (url models.ShortURL, err error) {
				return models.ShortURL{ID: "123456"}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			mockService := &mockURLShortenerService{
				CreateFunc: tt.mockFunc,
			}
			sh := &Handlers{
				service: mockService,
				baseURL: "http://example.com",
				logger:  logging.NewLogrusLogger(logrus.DebugLevel),
				auth:    auth.NewAuth(),
			}
			router.POST("/shorten", sh.HandleShorten)

			req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(tt.input))
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}
