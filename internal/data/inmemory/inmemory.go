package data

import (
	"errors"
	"fmt"
	"github.com/maxzhirnov/urlshort/internal/models"
)

type Store struct {
	urls map[string]models.URL
}

func NewInMemoryStore(urls map[string]models.URL) *Store {
	return &Store{
		urls: urls,
	}
}

func (s *Store) Save(url models.URL) error {
	if _, ok := s.urls[url.ID]; ok {
		return fmt.Errorf("id: %s already presented", url.ID)
	}

	s.urls[url.ID] = url
	return nil
}

func (s *Store) Get(id string) (models.URL, error) {
	url, ok := s.urls[id]
	if !ok {
		return models.URL{}, errors.New("id not found")
	}
	return url, nil
}
