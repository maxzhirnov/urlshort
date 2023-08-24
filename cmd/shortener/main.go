package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/maxzhirnov/urlshort/internal/app"
	"github.com/maxzhirnov/urlshort/internal/config"
	"github.com/maxzhirnov/urlshort/internal/handlers"
	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/middleware"
	"github.com/maxzhirnov/urlshort/internal/repository"
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

	cfg, err := config.NewFromFlags()
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Starting app",
		"server_addr", cfg.ServerAddr(),
		"base_url", cfg.BaseURL(),
		"file_storage_path", cfg.FileStoragePath())

	var storage app.Storage
	if cfg.FileStoragePath() == "" {
		storage = repository.NewMemoryStorage(logger)
	} else {
		storage, err = repository.NewFileStorage(cfg.FileStoragePath(), logger)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	defer storage.Close()

	service := app.NewURLShortener(storage)
	shortenerHandlers := handlers.NewShortenerHandlers(service, cfg.BaseURL())

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(middleware.Logging(logger))
	r.Use(middleware.GzipMiddleware())

	r.GET("/:ID", shortenerHandlers.HandleRedirect())
	r.POST("/", shortenerHandlers.HandleCreate())

	api := r.Group("/api")
	api.POST("/shorten", shortenerHandlers.HandleShorten())

	if err := r.Run(cfg.ServerAddr()); err != nil {
		logger.Error("Couldn't start server",
			"error", err,
		)
	}
}
