package main

import (
	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/middleware"
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
	logger, err := logging.NewZapSugared()
	if err != nil {
		panic(err)
	}

	cfg := config.NewDefaultConfig()
	cfg.Parse()

	logger.Info("Starting app",
		"server_addr", cfg.ServerAddr,
		"base_url", cfg.BaseURL)

	storage := memorystorage.New()
	urlShortenerService := app.NewURLShortener(storage)
	shortenerHandlers := handlers.NewShortenerHandlers(urlShortenerService, cfg.BaseURL)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Logging(logger))
	r.GET("/:ID", shortenerHandlers.HandleRedirect())
	r.POST("/", shortenerHandlers.HandleCreate())
	api := r.Group("/api")
	api.POST("/shorten", shortenerHandlers.HandleShorten)

	if err := r.Run(cfg.ServerAddr); err != nil {
		logger.Error("Couldn't start server",
			"error", err,
		)
	}
}
