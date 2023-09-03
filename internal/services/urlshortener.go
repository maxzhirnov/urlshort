package services

import (
	"context"
	"errors"
	"time"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

type repository interface {
	Insert(context.Context, models.ShortURL) error
	InsertMany(context.Context, []models.ShortURL) error
	Get(ctx context.Context, id string) (*models.ShortURL, error)
	Ping() error
}

type idGenerator interface {
	Generate() string
}

type URLShortener struct {
	Repo        repository
	IDGenerator idGenerator
	logger      logger
}

func NewURLShortener(repo repository, idGenerator idGenerator, logger logger) *URLShortener {
	return &URLShortener{
		Repo:        repo,
		IDGenerator: idGenerator,
		logger:      logger,
	}
}

func (us URLShortener) Create(originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	urlShorten := models.ShortURL{
		OriginalURL: originalURL,
		ID:          us.IDGenerator.Generate(),
	}
	if originalURL == "" {
		return "", errors.New("originalURL shouldn't be empty string")
	}
	if err := us.Repo.Insert(ctx, urlShorten); err != nil {
		return "", err
	}
	return urlShorten.ID, nil
}

func (us URLShortener) Get(id string) (*models.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if id == "" {
		return &models.ShortURL{}, errors.New("id shouldn't be empty string")
	}
	return us.Repo.Get(ctx, id)
}

func (us URLShortener) CreateBatch(urls []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if len(urls) == 0 {
		return []string{}, nil
	}

	shortenURLs := make([]models.ShortURL, len(urls))

	for i, url := range urls {
		shortenURLs[i] = models.ShortURL{
			OriginalURL: url,
			ID:          us.IDGenerator.Generate(),
		}
	}

	if err := us.Repo.InsertMany(ctx, shortenURLs); err != nil {
		us.logger.Error(err.Error())
		return nil, err
	}

	ids := make([]string, len(urls))
	for i, u := range shortenURLs {
		ids[i] = u.ID
	}
	return ids, nil
}

func (us URLShortener) Ping() error {
	return us.Repo.Ping()
}
