package main

import (
	"github.com/maxzhirnov/urlshort/internal/app"
	data "github.com/maxzhirnov/urlshort/internal/data/inmemory"
	"github.com/maxzhirnov/urlshort/internal/models"
)

func main() {
	urlData := make(map[string]models.URL)
	store := data.NewInMemoryStore(urlData)
	urlShortener := app.NewURLShortener(store)
	server := NewServer(urlShortener)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
