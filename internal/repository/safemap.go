package repository

import (
	"sync"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type safeMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func newSafeMap() *safeMap {
	return &safeMap{
		m: make(map[string]string),
	}
}

func (sm *safeMap) Load(id string) (urlObject models.ShortURL, ok bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	urlObject = models.ShortURL{}
	urlObject.OriginalURL, ok = sm.m[id]
	if ok {
		urlObject.ID = id
	}
	return urlObject, ok
}

func (sm *safeMap) Store(url models.ShortURL) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[url.ID] = url.OriginalURL
}
