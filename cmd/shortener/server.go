package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/internal/models"
	"log"
)

type URLShortenerService interface {
	Create(originalURL string) (id string, err error)
	Get(id string) (url models.URL, err error)
}

type Server struct {
	URLShortener URLShortenerService
}

func NewServer(urlShortener URLShortenerService) *Server {
	return &Server{
		URLShortener: urlShortener,
	}
}

func (s Server) Run() error {
	log.Println("Starting server...")

	r := gin.Default()
	r.GET("/:ID", handleRedirectToOriginal(s.URLShortener))
	r.POST("/", handleCreate(s.URLShortener))
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
	return nil
}
