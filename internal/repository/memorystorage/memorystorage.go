package memorystorage

import (
	"errors"
	"fmt"
	"github.com/maxzhirnov/urlshort/internal/logging"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type MemoryStorage struct {
	logger logging.Logger
	sm     *SafeMap
}

func New(logger logging.Logger) *MemoryStorage {
	return &MemoryStorage{
		sm:     NewSafeMap(),
		logger: logger,
	}
}

func (s *MemoryStorage) Save(url models.ShortURL) error {
	s.logger.Info("invoking Save() method in memory storage repo")
	// check if url with this id already exists
	if _, ok := s.sm.Load(url.ID); ok {
		err := errors.New(fmt.Sprintf("id: %s already presented", url.ID))
		s.logger.Error("error saving url in memory storage repo",
			"errorMessage", err)
		return err
	}

	s.sm.Store(url)
	s.logger.Info("Saved new url with memory storage repo",
		"url", url.OriginalURL,
		"id", url.ID)

	return nil
}

func (s *MemoryStorage) Get(id string) (models.ShortURL, error) {
	s.logger.Info("invoking Get() method in memory storage repo")
	url, ok := s.sm.Load(id)
	if !ok {
		err := errors.New("id not found")
		s.logger.Error("error getting url by id in memory storage repo",
			"errorMessage", err)
		return models.ShortURL{}, err
	}
	s.logger.Info("found url by id in memory storage repo",
		"url", url.OriginalURL,
		"id", url.ID)

	return url, nil
}

func (s *MemoryStorage) Exists(id string) bool {
	s.logger.Info("invoking Exists() method in memory storage repo")
	_, ok := s.sm.Load(id)
	s.logger.Info("Result of Exists()",
		"result", ok)
	return ok
}
