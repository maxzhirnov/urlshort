package storages

import (
	"context"

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

func (s *CombinedStorage) Insert(ctx context.Context, url models.ShortURL) (models.ShortURL, error) {
	if _, err := s.safeMap.Insert(context.Background(), url); err != nil {
		return models.ShortURL{}, err
	}

	if _, err := s.safeFile.Insert(context.Background(), url); err != nil {
		return models.ShortURL{}, err
	}

	return url, nil
}

func (s *CombinedStorage) InsertMany(ctx context.Context, urls []models.ShortURL) error {
	if err := s.safeMap.InsertMany(ctx, urls); err != nil {
		return err
	}

	if err := s.safeFile.InsertMany(ctx, urls); err != nil {
		return err
	}

	return nil
}

func (s *CombinedStorage) Get(ctx context.Context, id string) (models.ShortURL, bool) {
	return s.safeMap.Get(context.Background(), id)
}

func (s *CombinedStorage) Bootstrap(ctx context.Context) error {
	return s.safeFile.initializeData(s.safeMap)
}

func (s *CombinedStorage) Ping() error {
	return nil
}

func (s *CombinedStorage) Close() error {
	return s.safeFile.Close()
}
