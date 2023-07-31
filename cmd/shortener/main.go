package main

import (
	"github.com/maxzhirnov/urlshort/cmd/shortener/config"
	"github.com/maxzhirnov/urlshort/internal/app"
	data "github.com/maxzhirnov/urlshort/internal/data/inmemory"
)

func main() {
	cfg := config.NewDefaultConfig()
	cfg.ParseFlags()
	cfg.RedirectURLProtocol = config.HTTP

	safeMap := data.NewSafeMap()
	storeService := data.NewInMemoryStore(safeMap)
	urlShortenerService := app.NewURLShortener(storeService)
	server := NewServer(urlShortenerService, cfg)
	server.Run()
}
