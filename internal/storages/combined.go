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

func (s *CombinedStorage) InsertURL(ctx context.Context, url models.ShortURL) (models.ShortURL, error) {
	if _, err := s.safeMap.InsertURL(context.Background(), url); err != nil {
		return models.ShortURL{}, err
	}

	if _, err := s.safeFile.InsertURL(context.Background(), url); err != nil {
		return models.ShortURL{}, err
	}

	return url, nil
}

func (s *CombinedStorage) InsertURLMany(ctx context.Context, urls []models.ShortURL) error {
	if err := s.safeMap.InsertURLMany(ctx, urls); err != nil {
		return err
	}

	if err := s.safeFile.InsertURLMany(ctx, urls); err != nil {
		return err
	}

	return nil
}

func (s *CombinedStorage) GetURLByID(ctx context.Context, id string) (models.ShortURL, bool) {
	return s.safeMap.GetURLByID(context.Background(), id)
}

func (s *CombinedStorage) GetURLByOriginalURL(ctx context.Context, url string) (models.ShortURL, bool) {
	return s.safeMap.GetURLByOriginalURL(context.Background(), url)
}

func (s *CombinedStorage) GetURLsByUUID(ctx context.Context, uuid string) ([]models.ShortURL, error) {
	return s.safeMap.GetURLsByUUID(ctx, uuid)
}

func (s *CombinedStorage) Bootstrap(ctx context.Context) error {
	return s.safeFile.initializeData(s.safeMap)
}

func (s *CombinedStorage) TagURLsDeleted(ctx context.Context, urlsToDelete []models.Deletion) error {
	return nil
}

func (s *CombinedStorage) Ping() error {
	return nil
}

func (s *CombinedStorage) Close() error {
	return s.safeFile.Close()
}
