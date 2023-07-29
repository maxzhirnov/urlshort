package app

import (
	"github.com/maxzhirnov/urlshort/internal/models"
)

type StoreService interface {
	Save(url models.URL) error
	Get(id string) (models.URL, error)
}

type URLShortener struct {
	Store StoreService
}

func NewURLShortener(store StoreService) *URLShortener {
	return &URLShortener{
		Store: store,
	}
}

func (us URLShortener) Create(originalURL string) (id string, err error) {
	urlShorten := models.URL{
		OriginalURL: originalURL,
		ID:          generateID(8),
	}
	if err := us.Store.Save(urlShorten); err != nil {
		return "", err
	}
	return urlShorten.ID, nil
}

func (us URLShortener) Get(id string) (url models.URL, err error) {
	return us.Store.Get(id)
}
