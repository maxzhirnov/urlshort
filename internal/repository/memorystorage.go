package repository

import (
	"errors"
	"fmt"

	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/models"
)

// PersistentMemoryStorage save shortURL data in memory as well in file if Persistent()
// method called with fileStoragePath parameter
type PersistentMemoryStorage struct {
	logger logging.Logger
	memory *safeMap
	file   *safeFile
}

func NewPersistentMemoryStorage(logger logging.Logger) *PersistentMemoryStorage {
	return &PersistentMemoryStorage{
		logger: logger,
		memory: newSafeMap(),
		file:   nil,
	}
}

func (s *PersistentMemoryStorage) Create(url models.ShortURL) error {
	s.logger.Info("invoking Create() method in memory storage repo")
	// check if url with this id already exists
	if _, ok := s.memory.Load(url.ID); ok {
		err := fmt.Errorf("id: %s already presented", url.ID)
		s.logger.Error("error saving url in memory storage repo",
			"errorMessage", err)
		return err
	}

	s.memory.Store(url)
	s.logger.Info("Saved new url with memory storage repo",
		"url", url.OriginalURL,
		"id", url.ID)

	if err := s.file.Store(url); err != nil {
		return err
	}

	return nil
}

func (s *PersistentMemoryStorage) WithFileStorage(path string) error {
	if file, err := newSafeFile(path); err != nil {
		return err
	} else {
		s.file = file
	}
	return nil
}

func (s *PersistentMemoryStorage) LoadFileData() error {
	if s.file == nil {
		return fmt.Errorf("no file path was provided for PersostentMemoryStorage")
	}
	shortURLs, err := s.file.LoadAll()
	if err != nil {
		return err
	}

	for _, u := range shortURLs {
		s.memory.Store(u)
	}
	return nil
}

func (s *PersistentMemoryStorage) Get(id string) (models.ShortURL, error) {
	s.logger.Info("invoking Get() method in memory storage repo")
	url, ok := s.memory.Load(id)
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

func (s *PersistentMemoryStorage) Close() error {
	if s.file != nil {
		err := s.file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
