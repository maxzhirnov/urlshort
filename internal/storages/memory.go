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

func (s *MemoryStorage) GetURLByID(ctx context.Context, id string) (models.ShortURL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url := models.ShortURL{}
	var ok bool
	url.OriginalURL, ok = s.m[id]
	if ok {
		url.ID = id
	}
	return url, ok
}

func (s *MemoryStorage) GetURLByOriginalURL(ctx context.Context, url string) (models.ShortURL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var urlFound models.ShortURL
	for k, v := range s.m {
		if v == url {
			urlFound.OriginalURL = v
			urlFound.ID = k
			return urlFound, true
		}
	}
	return urlFound, false
}

func (s *MemoryStorage) InsertURL(ctx context.Context, url models.ShortURL) (models.ShortURL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[url.ID] = url.OriginalURL
	return url, nil
}

func (s *MemoryStorage) InsertURLMany(ctx context.Context, urls []models.ShortURL) error {
	for _, url := range urls {
		if _, err := s.InsertURL(ctx, url); err != nil {
			return err
		}
	}
	return nil
}

func (s *MemoryStorage) TagURLsDeleted(ctx context.Context, urlsToDelete []models.Deletion) error {
	return nil
}

func (s *MemoryStorage) GetURLsByUUID(ctx context.Context, uuid string) ([]models.ShortURL, error) {
	return nil, nil
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
