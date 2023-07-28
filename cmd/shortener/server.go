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
	URLShortener URLShortenerService
}

func NewServer(urlShortener URLShortenerService) *Server {
	return &Server{
		URLShortener: urlShortener,
	}
}

func (s Server) Run() error {
	log.Println("Starting server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//Разводим по методы тут, чтобы в хендлерах не пришлось объединять 2 хендлера в один
		//В дальнейшем с gin по идее можно будет сделать элегантнее
		if r.Method == http.MethodGet {
			h := handleGetOriginalURLByID(s.URLShortener)
			h(w, r)
		} else if r.Method == http.MethodPost {
			h := handleCreateShortURL(s.URLShortener)
			h(w, r)
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}

	return nil
}
