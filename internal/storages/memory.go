package storages

import (
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

func (ms *MemoryStorage) Get(id string) (urlObject *models.ShortURL, ok bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	urlObject = &models.ShortURL{}
	urlObject.OriginalURL, ok = ms.m[id]
	if ok {
		urlObject.ID = id
	}
	return urlObject, ok
}

func (ms *MemoryStorage) Insert(url models.ShortURL) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.m[url.ID] = url.OriginalURL
	return nil
}

func (ms *MemoryStorage) Ping() error {
	return nil
}
