package storages

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
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

func (fs *FileStorage) Insert(url models.ShortURL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	if _, err := fs.writer.Write(data); err != nil {
		return err
	}

	if err := fs.writer.WriteByte('\n'); err != nil {
		return err
	}

	err = fs.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) Get(id string) (*models.ShortURL, bool) {
	_, err := fs.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, false
	}
	fs.scanner = bufio.NewScanner(fs.file)
	for fs.scanner.Scan() {
		var shortURL *models.ShortURL
		err := json.Unmarshal(fs.scanner.Bytes(), &shortURL)
		if err != nil {
			continue
		}

		if shortURL.ID == id {
			return shortURL, true
		}
	}

	if err := fs.scanner.Err(); err != nil {
		return nil, false
	}

	return &models.ShortURL{}, false
}

func (fs *FileStorage) loadAll() ([]models.ShortURL, error) {
	shortURLs := make([]models.ShortURL, 0)
	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		var shortURL models.ShortURL
		err := json.Unmarshal(data, &shortURL)
		if err != nil {
			return nil, err
		}
		shortURLs = append(shortURLs, shortURL)
	}

	if err := fs.scanner.Err(); err != nil {
		return nil, err
	}

	return shortURLs, nil
}

// InitializeData loads all urls from file and upload them into memory storage
func (fs *FileStorage) InitializeData(memoryStorage *MemoryStorage) error {
	urls, err := fs.loadAll()
	if err != nil {
		return err
	}
	for _, u := range urls {
		if err := memoryStorage.Insert(u); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) Ping() error {
	return nil
}
