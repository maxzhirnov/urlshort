package data

import (
	"github.com/maxzhirnov/urlshort/internal/models"
	"sync"
)

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]string),
	}
}

func (sm *SafeMap) Load(key string) (urlObject models.URL, ok bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	url, ok := sm.m[key]
	urlObject = models.URL{
		OriginalURL: url,
		ID:          key,
	}
	return urlObject, ok
}

func (sm *SafeMap) Store(url models.URL) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[url.ID] = url.OriginalURL
}
