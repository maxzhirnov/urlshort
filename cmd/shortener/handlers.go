package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
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
		//_, isValid := app.CheckURL(originalHost)
		//if !isValid {
		//	//TODO: write test on that case
		//	c.String(http.StatusBadRequest, "provided data is not an URL")
		//}

		id, err := urlShortener.Create(originalHost)
		if err != nil {
			c.String(http.StatusInternalServerError, "error creating shorten url")
			return
		}
		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusCreated)
		shortenURL := fmt.Sprintf("%s/%s", redirectHost, id)

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
