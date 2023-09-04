package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/maxzhirnov/urlshort/internal/storages"
)

var ErrEntityAlreadyExist = errors.New("entity already exist")

type logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
}

// Storage is an interface for Storage services for storing and loading the url data
type Storage interface {
	Insert(context.Context, models.ShortURL) (models.ShortURL, error)
	InsertMany(context.Context, []models.ShortURL) error
	Get(ctx context.Context, id string) (models.ShortURL, bool)
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

func (r *Repository) Insert(ctx context.Context, url models.ShortURL) (models.ShortURL, error) {
	insertedURL, err := r.storage.Insert(ctx, url)
	if err != nil {
		if errors.Is(err, storages.ErrEntityAlreadyExist) {
			return insertedURL, ErrEntityAlreadyExist
		}
		r.logger.Error("error: ", err)
		return models.ShortURL{}, err
	}
	r.logger.Debug("Inserted successfully: ", url)
	return insertedURL, nil
}

func (r *Repository) InsertMany(ctx context.Context, urls []models.ShortURL) error {
	return r.storage.InsertMany(ctx, urls)
}

func (r *Repository) Get(ctx context.Context, id string) (models.ShortURL, error) {
	url, ok := r.storage.Get(ctx, id)
	if !ok {
		return models.ShortURL{}, fmt.Errorf("id not found")
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
