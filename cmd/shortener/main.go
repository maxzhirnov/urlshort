package main

import (
	"github.com/joho/godotenv"
	"github.com/maxzhirnov/urlshort/cmd/shortener/config"
	"github.com/maxzhirnov/urlshort/internal/app"
	data "github.com/maxzhirnov/urlshort/internal/data/inmemory"
	"log"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file parsing failed")
	}
}

func main() {
	cfg := config.NewDefaultConfig()
	cfg.Parse()

	log.Println("Starting app with config:")
	log.Printf("Server host: %s; Redirect address: %s", cfg.ServerAddr, cfg.BaseURL)

	safeMap := data.NewSafeMap()
	storeService := data.NewInMemoryStore(safeMap)
	urlShortenerService := app.NewURLShortener(storeService)
	server := NewServer(urlShortenerService, cfg)
	server.Run()
}
