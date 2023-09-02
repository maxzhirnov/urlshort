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

func (s *CombinedStorage) Insert(url models.ShortURL) error {
	if err := s.safeMap.Insert(url); err != nil {
		return err
	}

	if err := s.safeFile.Insert(url); err != nil {
		return err
	}

	return nil
}

func (s *CombinedStorage) Get(id string) (*models.ShortURL, bool) {
	return s.safeMap.Get(id)
}

func (s *CombinedStorage) Ping() error {
	return nil
}

func (s *CombinedStorage) Close() error {
	return s.safeFile.Close()
}
