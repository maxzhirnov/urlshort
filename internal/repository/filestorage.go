package repository

import (
	"github.com/maxzhirnov/urlshort/internal/logging"
	"github.com/maxzhirnov/urlshort/internal/models"
)

type FileStorage struct {
	filePath      string
	safeFile      *safeFile
	memoryStorage *MemoryStorage
	logger        logging.Logger
}

func NewFileStorage(filePath string, logger logging.Logger) (*FileStorage, error) {
	ms := NewMemoryStorage(logger)
	safeFile, err := newSafeFile(filePath)
	if err != nil {
		return nil, err
	}

	shortURLs, err := safeFile.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, u := range shortURLs {
		ms.sm.Store(u)
	}

	return &FileStorage{
		filePath:      filePath,
		safeFile:      safeFile,
		logger:        logger,
		memoryStorage: ms,
	}, nil
}

func (s *FileStorage) Create(url models.ShortURL) error {
	if err := s.safeFile.Store(url); err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return s.memoryStorage.Create(url)
}

func (s *FileStorage) Get(id string) (models.ShortURL, error) {
	return s.memoryStorage.Get(id)
}

func (s *FileStorage) Close() error {
	return s.safeFile.Close()
}
