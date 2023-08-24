package app

import (
	"errors"
	"github.com/maxzhirnov/urlshort/internal/models"
)

type Storage interface {
	Create(url models.ShortURL) error
	Get(id string) (models.ShortURL, error)
	Close() error
}

type URLShortener struct {
	Storage Storage
}

func NewURLShortener(s Storage) *URLShortener {
	return &URLShortener{
		Storage: s,
	}
}

func (us URLShortener) Create(originalURL string) (string, error) {
	urlShorten := models.ShortURL{
		OriginalURL: originalURL,
		ID:          generateID(8),
	}
	if originalURL == "" {
		return "", errors.New("originalURL shouldn't be empty string")
	}
	if err := us.Storage.Create(urlShorten); err != nil {
		return "", err
	}
	return urlShorten.ID, nil
}

func (us URLShortener) Get(id string) (models.ShortURL, error) {
	if id == "" {
		return models.ShortURL{}, errors.New("id shouldn't be empty string")
	}
	return us.Storage.Get(id)
}
