package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/app"
	"io"
	"net/http"
)

func handleCreate(urlShortener URLShortenerService, redirectHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.String(http.StatusBadRequest, "only POST requests allowed")
			return
		}

		defer c.Request.Body.Close()

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading request body")
			return
		}

		originalHost := string(data)

		id, err := urlShortener.Create(originalHost)
		if err != nil {
			c.String(http.StatusInternalServerError, "error creating shorten url")
			return
		}

		shortenURL := fmt.Sprintf("%s/%s", redirectHost, id)

		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusCreated)
		//Отдаем в body ссылку на сокращенный url
		if _, err := c.Writer.Write([]byte(shortenURL)); err != nil {
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}
	}
}

func handleRedirect(urlShortener URLShortenerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.String(http.StatusBadRequest, "only GET requests allowed")
			return
		}

		id := c.Param("ID")
		url, err := urlShortener.Get(id)
		if err != nil {
			c.String(http.StatusNotFound, "id not found")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, app.EnsureURLScheme(url.OriginalURL))
	}
}
