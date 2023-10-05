package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	InsertURL(context.Context, models.ShortURL) (models.ShortURL, error)
	InsertURLMany(context.Context, []models.ShortURL) error
	GetURLByID(ctx context.Context, id string) (models.ShortURL, bool)
	GetURLByOriginalURL(ctx context.Context, url string) (models.ShortURL, bool)
	GetURLsByUUID(ctx context.Context, uuid string) ([]models.ShortURL, error)
	TagURLsDeleted(context.Context, []models.Deletion) error
	Bootstrap() error
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
	insertedURL, err := r.storage.InsertURL(ctx, url)
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

// InsertMany inserts urls if not exists, if exists returns existing url object
func (r *Repository) InsertMany(ctx context.Context, urlsToInsert []models.ShortURL) ([]models.ShortURL, error) {
	if err := r.storage.InsertURLMany(ctx, urlsToInsert); err != nil {
		return nil, err
	}

	shortenURLs := make([]models.ShortURL, len(urlsToInsert))
	for i, u := range urlsToInsert {
		existingURL, _ := r.storage.GetURLByOriginalURL(ctx, u.OriginalURL)
		shortenURLs[i] = existingURL
	}

	return shortenURLs, nil
}

func (r *Repository) GetURLByID(ctx context.Context, id string) (models.ShortURL, error) {
	url, ok := r.storage.GetURLByID(ctx, id)
	if !ok {
		return models.ShortURL{}, fmt.Errorf("id not found")
	}
	return url, nil
}

func (r *Repository) GetURLsByUUID(ctx context.Context, uuid string) ([]models.ShortURL, error) {
	return r.storage.GetURLsByUUID(ctx, uuid)
}

func (r *Repository) TagURLsDeleted(urlsToDelete []models.Deletion) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.storage.TagURLsDeleted(ctx, urlsToDelete)
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
