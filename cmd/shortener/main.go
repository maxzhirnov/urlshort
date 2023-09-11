package main

import (
	"compress/gzip"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/maxzhirnov/urlshort/internal/configs"
	"github.com/maxzhirnov/urlshort/internal/handlers"
	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/middleware"
	"github.com/maxzhirnov/urlshort/internal/repositories"
	"github.com/maxzhirnov/urlshort/internal/services"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file parsing failed")
	}

	logger := logging.NewLogrusLogger(logrus.DebugLevel)

	config, err := configs.NewFromFlags(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Starting app",
		"server_addr", config.ServerAddr(),
		"base_url", config.BaseURL(),
		"file_storage_path", config.FileStoragePath())

	storage, err := repositories.NewStorage(*config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer storage.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := storage.Bootstrap(ctx); err != nil {
		logger.Fatal(err.Error())
	}

	repo := repositories.NewRepository(logger, storage)
	idGenerator := services.NewRandIDGenerator(8)
	service := services.NewURLShortener(repo, idGenerator, logger)
	handler := handlers.NewHandlers(service, config.BaseURL(), logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Logging(logger))
	gzipWriter, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
	if err != nil {
		log.Fatal(err)
	}
	r.Use(middleware.Gzip(logger, gzipWriter))

	r.GET("/:ID", handler.HandleRedirect)
	r.POST("/", handler.HandleCreate)
	r.GET("/ping", handler.HandlePing)

	api := r.Group("/api")
	api.POST("/shorten", handler.HandleShorten)
	api.POST("/shorten/batch", handler.HandleShortenBatch)

	if err := r.Run(config.ServerAddr()); err != nil {
		logger.Fatal("Couldn't start server",
			"error", err,
		)
	}
}
