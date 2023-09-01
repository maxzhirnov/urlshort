package main

import (
	"compress/gzip"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/maxzhirnov/urlshort/internal/configs"
	"github.com/maxzhirnov/urlshort/internal/handlers"
	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/middleware"
	"github.com/maxzhirnov/urlshort/internal/repository"
	"github.com/maxzhirnov/urlshort/internal/services"
	"github.com/maxzhirnov/urlshort/internal/storages"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file parsing failed")
	}
}

func main() {
	logger, err := logging.NewZapSugared()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	config, err := configs.NewFromFlags(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Starting services",
		"server_addr", config.ServerAddr(),
		"base_url", config.BaseURL(),
		"file_storage_path", config.FileStoragePath())

	var storage repository.Storage
	memoryStorage := storages.NewMemoryStorage()
	if config.ShouldSaveToFile() {
		fileStorage, err := storages.NewFileStorage(config.FileStoragePath())
		if err != nil {
			logger.Fatal(err.Error())
		}

		if err := fileStorage.InitializeData(memoryStorage); err != nil {
			logger.Fatal(err.Error())
		}

		storage = storages.NewCombinedStorage(memoryStorage, fileStorage)
	} else {
		storage = memoryStorage
	}

	postgresDB, err := storages.NewPostgresql(config.PostgresConn())
	if err != nil {
		logger.Fatal(err.Error())
	}

	repo := repository.NewRepository(logger, storage)
	postgresRepo := repository.NewRepository(logger, postgresDB)

	idGenerator := services.NewRandIDGenerator(8)
	service := services.NewURLShortener(repo, idGenerator)
	serviceWithPostgres := services.NewURLShortener(postgresRepo, idGenerator)

	shortenerHandlers := handlers.NewShortenerHandlers(service, config.BaseURL())
	handlersWithPostgres := handlers.NewShortenerHandlers(serviceWithPostgres, config.BaseURL())

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(middleware.Logging(logger))

	gzipWriter, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
	if err != nil {
		log.Fatal(err)
	}
	r.Use(middleware.Gzip(logger, gzipWriter))

	r.GET("/:ID", shortenerHandlers.HandleRedirect)
	r.POST("/", shortenerHandlers.HandleCreate)
	r.GET("/ping", handlersWithPostgres.HandlePing)

	api := r.Group("/api")
	api.POST("/shorten", shortenerHandlers.HandleShorten)

	if err := r.Run(config.ServerAddr()); err != nil {
		logger.Error("Couldn't start server",
			"error", err,
		)
	}
}
