package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/maxzhirnov/urlshort/internal/services"
)

type Service interface {
	Create(originalURL string) (id string, err error)
	Get(id string) (url *models.ShortURL, err error)
	Ping() error
}

type ShortenerHandlers struct {
	service Service
	baseURL string
}

func NewShortenerHandlers(s Service, baseURL string) *ShortenerHandlers {
	return &ShortenerHandlers{
		service: s,
		baseURL: baseURL,
	}
}

func (sh *ShortenerHandlers) HandleCreate(c *gin.Context) {
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

func (sh *ShortenerHandlers) HandleRedirect(c *gin.Context) {
	id := c.Param("ID")
	url, err := sh.service.Get(id)
	if err != nil {
		c.String(http.StatusNotFound, "id not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, services.EnsureURLScheme(url.OriginalURL))
}

func (sh *ShortenerHandlers) HandleShorten(c *gin.Context) {
	var reqData models.ShortenRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you should provide correct data"})
		return
	}
	defer c.Request.Body.Close()

	if len(reqData.URL) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url should be valid url"})
		return
	}

	shortenID, err := sh.service.Create(reqData.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	shortenURL := sh.baseURL + "/" + shortenID
	response := models.ShortenResponse{Result: shortenURL}
	c.JSON(http.StatusCreated, response)
}

func (sh *ShortenerHandlers) HandlePing(c *gin.Context) {
	err := sh.service.Ping()
	if err != nil {
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}
	c.String(http.StatusOK, "connected to database")
}
