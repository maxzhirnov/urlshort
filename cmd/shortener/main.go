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

	memoryStorage := repository.NewPersistentMemoryStorage(logger)

	if cfg.FileStoragePath() != "" {
		if err := memoryStorage.WithFileStorage(cfg.FileStoragePath()); err != nil {
			logger.Error(err.Error())
		}
		if err := memoryStorage.LoadFileData(); err != nil {
			logger.Error(err.Error())
		}
		defer memoryStorage.Close()
	}

	service := app.NewURLShortener(memoryStorage)
	shortenerHandlers := handlers.NewShortenerHandlers(service, cfg.BaseURL())

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(middleware.Logging(logger))
	r.Use(middleware.Gzip())

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
