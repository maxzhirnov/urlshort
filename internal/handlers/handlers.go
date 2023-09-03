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

type logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

type Service interface {
	Create(originalURL string) (id string, err error)
	CreateBatch([]string) (ids []string, err error)
	Get(id string) (url *models.ShortURL, err error)
	Ping() error
}

type Handlers struct {
	service Service
	baseURL string
	logger  logger
}

func NewHandlers(s Service, baseURL string, logger logger) *Handlers {
	return &Handlers{
		service: s,
		baseURL: baseURL,
		logger:  logger,
	}
}

func (h *Handlers) HandleCreate(c *gin.Context) {
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

	id, err := h.service.Create(originalHost)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating shorten url")
		return
	}

	shortenURL := fmt.Sprintf("%s/%s", h.baseURL, id)

	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(http.StatusCreated)
	//Отдаем в body ссылку на сокращенный url
	if _, err := c.Writer.Write([]byte(shortenURL)); err != nil {
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}
}

func (h *Handlers) HandleRedirect(c *gin.Context) {
	id := c.Param("ID")
	url, err := h.service.Get(id)
	if err != nil {
		c.String(http.StatusNotFound, "id not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, services.EnsureURLScheme(url.OriginalURL))
}

func (h *Handlers) HandleShorten(c *gin.Context) {
	var reqData struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you should provide correct data"})
		return
	}
	defer c.Request.Body.Close()

	if len(reqData.URL) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url should be valid url"})
		return
	}

	shortenID, err := h.service.Create(reqData.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	shortenURL := h.baseURL + "/" + shortenID

	response := struct {
		Result string `json:"result"`
	}{
		Result: shortenURL,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) HandleShortenBatch(c *gin.Context) {
	var request = make([]struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}, 0)

	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	defer c.Request.Body.Close()

	urlsToShort := make([]string, 0)
	for _, u := range request {
		urlsToShort = append(urlsToShort, u.OriginalURL)
	}

	ids, err := h.service.CreateBatch(urlsToShort)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	response := make([]struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}, len(ids))

	for i, id := range ids {
		response[i].CorrelationID = request[i].CorrelationID
		response[i].ShortURL = h.baseURL + "/" + id
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) HandlePing(c *gin.Context) {
	err := h.service.Ping()
	if err != nil {
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}
	c.String(http.StatusOK, "connected to database")
}
