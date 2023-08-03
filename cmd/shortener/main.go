package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/maxzhirnov/urlshort/internal/app"
	"github.com/maxzhirnov/urlshort/internal/config"
	"github.com/maxzhirnov/urlshort/internal/handlers"
	"github.com/maxzhirnov/urlshort/internal/repository/memorystorage"
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

	storage := memorystorage.New()
	urlShortenerService := app.NewURLShortener(storage)
	shortenerHandlers := handlers.NewShortenerHandlers(urlShortenerService, cfg.BaseURL)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:ID", shortenerHandlers.HandleRedirect())
	r.POST("/", shortenerHandlers.HandleCreate())

	if err := r.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("Couldn't start server: %s:", err)
	}
}
