package handlers

import (
	"fmt"
	"github.com/maxzhirnov/urlshort/internal/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/app"
)

type urlShortenerService interface {
	Create(originalURL string) (id string, err error)
	Get(id string) (url models.ShortURL, err error)
}

type ShortenerHandlers struct {
	service urlShortenerService
	baseURL string
}

func NewShortenerHandlers(s urlShortenerService, baseURL string) *ShortenerHandlers {
	return &ShortenerHandlers{
		service: s,
		baseURL: baseURL,
	}
}

func (sh *ShortenerHandlers) HandleCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading request body")
			return
		}

		if len(data) == 0 {
			c.String(http.StatusBadRequest, "url shouldn't be empty")
			return
		}

		originalHost := string(data)

		id, err := sh.service.Create(originalHost)
		if err != nil {
			c.String(http.StatusInternalServerError, "error creating shorten url")
			return
		}

		shortenURL := fmt.Sprintf("%s/%s", sh.baseURL, id)

		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusCreated)
		//Отдаем в body ссылку на сокращенный url
		if _, err := c.Writer.Write([]byte(shortenURL)); err != nil {
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}
	}
}

func (sh *ShortenerHandlers) HandleRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("ID")
		url, err := sh.service.Get(id)
		if err != nil {
			c.String(http.StatusNotFound, "id not found")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, app.EnsureURLScheme(url.OriginalURL))
	}
}
