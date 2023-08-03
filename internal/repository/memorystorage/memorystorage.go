package memorystorage

import (
	"errors"
	"fmt"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type MemoryStorage struct {
	sm *SafeMap
}

func New() *MemoryStorage {
	return &MemoryStorage{
		sm: NewSafeMap(),
	}
}

func (s *MemoryStorage) Save(url models.URL) error {
	// check if url with this id already exists
	if _, ok := s.sm.Load(url.ID); ok {
		return fmt.Errorf("id: %s already presented", url.ID)
	}

	s.sm.Store(url)
	return nil
}

func (s *MemoryStorage) Get(id string) (models.URL, error) {
	url, ok := s.sm.Load(id)
	if !ok {
		return models.URL{}, errors.New("id not found")
	}
	return url, nil
}
