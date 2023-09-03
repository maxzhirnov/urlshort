package repository

import (
	"context"
	"fmt"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

// Storage is an interface for Storage services for storing and loading the url data
type Storage interface {
	Insert(context.Context, models.ShortURL) error
	InsertMany(context.Context, []models.ShortURL) error
	Get(ctx context.Context, id string) (*models.ShortURL, bool)
	Bootstrap(context.Context) error
	Close() error
	Ping() error
}

// Repository save shortURL data in storages as well in persistentStorage if WithFileStorage()
// method called with fileStoragePath parameter
type Repository struct {
	logger  logger
	storage Storage
}

func NewRepository(logger logger, storage Storage) *Repository {
	return &Repository{
		logger:  logger,
		storage: storage,
	}
}

func (r *Repository) Insert(ctx context.Context, url models.ShortURL) error {
	// check for id collision
	if _, ok := r.storage.Get(ctx, url.ID); ok {
		r.logger.Error("id collision while trying to crate new short url", "id", url.ID)
		return fmt.Errorf("please try again")
	}
	return r.storage.Insert(ctx, url)
}

func (r *Repository) InsertMany(ctx context.Context, urls []models.ShortURL) error {
	return r.storage.InsertMany(ctx, urls)
}

func (r *Repository) Get(ctx context.Context, id string) (*models.ShortURL, error) {
	url, ok := r.storage.Get(ctx, id)
	if !ok {
		return nil, fmt.Errorf("id not found")
	}
	return url, nil
}

func (r *Repository) Ping() error {
	err := r.storage.Ping()
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	return nil
}

func (r *Repository) Close() error {
	return r.storage.Close()
}
