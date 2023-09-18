package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxzhirnov/urlshort/internal/auth"
	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/maxzhirnov/urlshort/internal/services"
)

type logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
}

type service interface {
	Create(url, uuid string) (models.ShortURL, error)
	CreateBatch(urls []string, uuid string) (ids []string, err error)
	Get(id string) (url models.ShortURL, err error)
	GetAllUsersURLs(uuid string) ([]models.ShortURL, error)
	Ping() error
}

type Handlers struct {
	service service
	baseURL string
	auth    *auth.Auth
	logger  logger
}

func NewHandlers(s service, baseURL string, auth *auth.Auth, logger logger) *Handlers {
	return &Handlers{
		service: s,
		baseURL: baseURL,
		auth:    auth,
		logger:  logger,
	}
}

func (h *Handlers) HandleCreate(c *gin.Context) {
	defer c.Request.Body.Close()
	originalURLData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading request body")
		return
	}

	if len(originalURLData) == 0 {
		c.String(http.StatusBadRequest, "url shouldn't be empty")
		return
	}

	originalURL := string(originalURLData)

	userID, err := h.getUserIDFromJWTToken(c)
	if err != nil {
		h.logger.Warn(err.Error())
	}

	statusCode := http.StatusCreated
	shortenURLObject, err := h.service.Create(originalURL, userID)

	if errors.Is(err, services.ErrEntityAlreadyExist) {
		statusCode = http.StatusConflict
	} else if err != nil {
		c.String(http.StatusInternalServerError, "error creating shorten url")
		return
	}

	shortenURL := fmt.Sprintf("%s/%s", h.baseURL, shortenURLObject.ID)

	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(statusCode)
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

	userID, err := h.getUserIDFromJWTToken(c)
	if err != nil {
		h.logger.Warn(err.Error())
	}

	statusCode := http.StatusCreated
	shortenURLObject, err := h.service.Create(reqData.URL, userID)
	if err != nil {
		switch {
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		case errors.Is(err, services.ErrEntityAlreadyExist):
			statusCode = http.StatusConflict
		}
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: h.baseURL + "/" + shortenURLObject.ID,
	}
	c.JSON(statusCode, response)
}

func (h *Handlers) HandleShortenBatch(c *gin.Context) {
	var request = make([]struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}, 0)

	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		h.logger.Error("error decoding json", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	defer c.Request.Body.Close()

	userID, err := h.getUserIDFromJWTToken(c)
	if err != nil {
		h.logger.Warn(err.Error())
	}

	urlsToShort := make([]string, 0)
	for _, u := range request {
		urlsToShort = append(urlsToShort, u.OriginalURL)
	}

	ids, err := h.service.CreateBatch(urlsToShort, userID)
	if err != nil {
		h.logger.Error("error creating batch", err)
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

type ShowAllUsersURLsDTO struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *Handlers) newShowAllUsersURLsDTO(su models.ShortURL) ShowAllUsersURLsDTO {
	return ShowAllUsersURLsDTO{
		ShortURL:    h.baseURL + "/" + su.ID,
		OriginalURL: "http://" + su.OriginalURL,
	}
}

func (h *Handlers) HandleShowAllUsersURLs(c *gin.Context) {
	jwtToken, err := c.Cookie("jwt_token")
	if err != nil {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := h.auth.ValidateToken(jwtToken)
	if err != nil {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	userURLs, err := h.service.GetAllUsersURLs(userID)
	if len(userURLs) == 0 {
		c.JSON(http.StatusNoContent, "empty")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]ShowAllUsersURLsDTO, len(userURLs))
	for i, u := range userURLs {
		dto := h.newShowAllUsersURLsDTO(u)
		res[i] = dto
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handlers) getUserIDFromJWTToken(c *gin.Context) (string, error) {
	var jwtToken string
	var err error
	// Пытаемся получить jwtToken из контекста
	if tempToken, exists := c.Get("jwt_token"); exists {
		jwtToken = tempToken.(string)
	} else {
		// Если токена нет в контексте, пытаемся получить его из куки
		jwtToken, err = c.Cookie("jwt_token")
		if err != nil {
			return "", err
		}
	}

	var userID string
	userID, err = h.auth.ValidateToken(jwtToken)
	if err != nil {
		return "", err
	}

	return userID, nil
}
