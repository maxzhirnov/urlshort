package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maxzhirnov/urlshort/cmd/shortener/config"
	"github.com/maxzhirnov/urlshort/internal/models"
	"log"
)

type URLShortenerService interface {
	Create(originalURL string) (id string, err error)
	Get(id string) (url models.URL, err error)
}

type Server struct {
	URLShortener        URLShortenerService
	ServerAddr          string
	RedirectHost        string
	RedirectURLProtocol string
}

func NewServer(urlShortener URLShortenerService, cfg *config.Config) *Server {
	return &Server{
		URLShortener:        urlShortener,
		ServerAddr:          cfg.ServerAddr,
		RedirectHost:        cfg.RedirectHost,
		RedirectURLProtocol: string(cfg.RedirectURLProtocol),
	}
}

func (s Server) Run() {
	log.Println("Starting server...")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:ID", handleRedirect(s.URLShortener))
	r.POST("/", handleCreate(s.URLShortener, s.RedirectURLProtocol+s.RedirectHost))

	if err := r.Run(s.ServerAddr); err != nil {
		panic(err)
	}
}
