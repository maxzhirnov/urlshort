package storages

import (
	"github.com/maxzhirnov/urlshort/internal/models"
)

type CombinedStorage struct {
	safeMap  *MemoryStorage
	safeFile *FileStorage
}

func NewCombinedStorage(safeMap *MemoryStorage, safeFile *FileStorage) *CombinedStorage {
	return &CombinedStorage{
		safeMap:  safeMap,
		safeFile: safeFile,
	}
}

func (s *CombinedStorage) Store(url models.ShortURL) error {
	if err := s.safeMap.Store(url); err != nil {
		return err
	}

	if err := s.safeFile.Store(url); err != nil {
		return err
	}

	return nil
}

func (s *CombinedStorage) Load(id string) (*models.ShortURL, bool) {
	return s.safeMap.Load(id)
}
