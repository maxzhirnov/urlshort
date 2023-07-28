package data

import (
	"github.com/maxzhirnov/urlshort/internal/models"
	"sync"
)

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]models.URL
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]models.URL),
	}
}

func (sm *SafeMap) Load(key string) (url models.URL, ok bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	url, ok = sm.m[key]
	return
}

func (sm *SafeMap) Store(url models.URL) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[url.ID] = url
}
