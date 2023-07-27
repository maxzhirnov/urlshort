package main

import (
	"github.com/maxzhirnov/urlshort/internal/models"
	"log"
	"net/http"
)

type URLShortenerService interface {
	Create(originalURL string) (id string, err error)
	Get(id string) (url models.URL, err error)
}

type Server struct {
	UrlShortener URLShortenerService
}

func NewServer(urlShortener URLShortenerService) *Server {
	return &Server{
		UrlShortener: urlShortener,
	}
}

func (s Server) Run() error {
	log.Println("Starting server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			f := handleGetOriginalURLByID(s.UrlShortener)
			f(w, r)
		} else if r.Method == http.MethodPost {
			f := handleCreateShortURL(s.UrlShortener)
			f(w, r)
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}

	return nil
}
