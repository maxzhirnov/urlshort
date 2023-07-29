package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/app"
	"io"
	"net/http"
	"strings"
)

func handleCreate(urlShortener URLShortenerService) gin.HandlerFunc {
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

		isValid, url := app.CheckURL(string(data))
		if !isValid {
			//TODO: write test on that case
			c.String(http.StatusBadRequest, "provided data is not an URL")
		}

		p := "http://" //TODO: Implement protocol parsing and mapping to string
		h := c.Request.Host
		id, err := urlShortener.Create(url)
		if err != nil {
			c.String(http.StatusInternalServerError, "error creating shorten url")
			return
		}
		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusCreated)
		shortenURL := fmt.Sprintf("%s%s/%s", p, h, id)

		if _, err := c.Writer.Write([]byte(shortenURL)); err != nil {
			c.String(http.StatusInternalServerError, "something went wrong")
			return
		}
	}
}

func handleRedirectToOriginal(urlShortener URLShortenerService) gin.HandlerFunc {
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

		originalURL := url.OriginalURL
		if !strings.HasPrefix(originalURL, "http") {
			originalURL = "https://" + originalURL
		}

		c.Redirect(http.StatusTemporaryRedirect, originalURL)
	}
}
