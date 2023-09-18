package storages

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
	mu      sync.RWMutex
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	if filePath == "" {
		return nil, nil
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (s *FileStorage) InsertURL(ctx context.Context, url models.ShortURL) (models.ShortURL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.Marshal(url)
	if err != nil {
		return models.ShortURL{}, err
	}

	if _, err := s.writer.Write(data); err != nil {
		return models.ShortURL{}, err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return models.ShortURL{}, err
	}

	err = s.writer.Flush()
	if err != nil {
		return models.ShortURL{}, err
	}
	return url, nil
}

func (s *FileStorage) InsertURLMany(ctx context.Context, urls []models.ShortURL) error {
	for _, url := range urls {
		if _, err := s.InsertURL(ctx, url); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) GetURLByID(ctx context.Context, id string) (models.ShortURL, bool) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return models.ShortURL{}, false
	}
	s.scanner = bufio.NewScanner(s.file)
	for s.scanner.Scan() {
		var shortURL models.ShortURL
		err := json.Unmarshal(s.scanner.Bytes(), &shortURL)
		if err != nil {
			continue
		}

		if shortURL.ID == id {
			return shortURL, true
		}
	}

	if err := s.scanner.Err(); err != nil {
		return models.ShortURL{}, false
	}

	return models.ShortURL{}, false
}

func (s *FileStorage) GetURLByOriginalURL(ctx context.Context, url string) (models.ShortURL, bool) {
	panic(fmt.Errorf("Not implemented"))
}

func (s *FileStorage) GetURLsByUUID(ctx context.Context, uuid string) ([]models.ShortURL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *FileStorage) Bootstrap(ctx context.Context) error {
	return nil
}

func (s *FileStorage) Ping() error {
	return nil
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}

// initializeData loads all urls from file and upload them into memory storage
func (s *FileStorage) initializeData(memoryStorage *MemoryStorage) error {
	urls, err := s.loadAll()
	if err != nil {
		return err
	}
	for _, u := range urls {
		if _, err := memoryStorage.InsertURL(context.Background(), u); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) loadAll() ([]models.ShortURL, error) {
	shortURLs := make([]models.ShortURL, 0)
	for s.scanner.Scan() {
		data := s.scanner.Bytes()
		var shortURL models.ShortURL
		err := json.Unmarshal(data, &shortURL)
		if err != nil {
			return nil, err
		}
		shortURLs = append(shortURLs, shortURL)
	}

	if err := s.scanner.Err(); err != nil {
		return nil, err
	}

	return shortURLs, nil
}
