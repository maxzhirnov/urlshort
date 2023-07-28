package data

import (
	"errors"
	"fmt"
	"github.com/maxzhirnov/urlshort/internal/models"
)

type Store struct {
	urls *SafeMap
}

func NewInMemoryStore(m *SafeMap) *Store {
	return &Store{
		urls: m,
	}
}

func (s *Store) Save(url models.URL) error {
	if _, ok := s.urls.Load(url.ID); ok {
		return fmt.Errorf("id: %s already presented", url.ID)
	}

	s.urls.Store(url)
	return nil
}

func (s *Store) Get(id string) (models.URL, error) {
	url, ok := s.urls.Load(id)
	if !ok {
		return models.URL{}, errors.New("id not found")
	}
	return url, nil
}
