package app

import (
	"errors"
	"github.com/maxzhirnov/urlshort/internal/models"
)

type storage interface {
	Save(url models.URL) error
	Get(id string) (models.URL, error)
}

type URLShortener struct {
	Storage storage
}

func NewURLShortener(s storage) *URLShortener {
	return &URLShortener{
		Storage: s,
	}
}

func (us URLShortener) Create(originalURL string) (string, error) {
	urlShorten := models.URL{
		OriginalURL: originalURL,
		ID:          generateID(8),
	}
	if originalURL == "" {
		return "", errors.New("originalURL shouldn't be empty string")
	}
	if err := us.Storage.Save(urlShorten); err != nil {
		return "", err
	}
	return urlShorten.ID, nil
}

func (us URLShortener) Get(id string) (models.URL, error) {
	if id == "" {
		return models.URL{}, errors.New("id shouldn't be empty string")
	}
	return us.Storage.Get(id)
}
