package services

import (
	"errors"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type repository interface {
	Create(url models.ShortURL) error
	Get(id string) (*models.ShortURL, error)
	Ping() error
}

type idGenerator interface {
	Generate() string
}

type URLShortener struct {
	Storage     repository
	IDGenerator idGenerator
}

func NewURLShortener(repo repository, idGenerator idGenerator) *URLShortener {
	return &URLShortener{
		Storage:     repo,
		IDGenerator: idGenerator,
	}
}

func (us URLShortener) Create(originalURL string) (string, error) {
	urlShorten := models.ShortURL{
		OriginalURL: originalURL,
		ID:          us.IDGenerator.Generate(),
	}
	if originalURL == "" {
		return "", errors.New("originalURL shouldn't be empty string")
	}
	if err := us.Storage.Create(urlShorten); err != nil {
		return "", err
	}
	return urlShorten.ID, nil
}

func (us URLShortener) Get(id string) (*models.ShortURL, error) {
	if id == "" {
		return &models.ShortURL{}, errors.New("id shouldn't be empty string")
	}
	return us.Storage.Get(id)
}

func (us URLShortener) Ping() error {
	return us.Storage.Ping()
}
