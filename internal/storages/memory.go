package storages

import (
	"context"
	"sync"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type MemoryStorage struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		m: make(map[string]string),
	}
}

func (s *MemoryStorage) Get(ctx context.Context, id string) (urlObject *models.ShortURL, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	urlObject = &models.ShortURL{}
	urlObject.OriginalURL, ok = s.m[id]
	if ok {
		urlObject.ID = id
	}
	return urlObject, ok
}

func (s *MemoryStorage) Insert(ctx context.Context, url models.ShortURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[url.ID] = url.OriginalURL
	return nil
}

func (s *MemoryStorage) InsertMany(ctx context.Context, urls []models.ShortURL) error {
	for _, url := range urls {
		if err := s.Insert(ctx, url); err != nil {
			return err
		}
	}
	return nil
}

func (s *MemoryStorage) Bootstrap(ctx context.Context) error {
	return nil
}

func (s *MemoryStorage) Ping() error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
