package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type safeFile struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func newSafeFile(filePath string) (*safeFile, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &safeFile{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (sf *safeFile) Store(url models.ShortURL) error {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	if _, err := sf.writer.Write(data); err != nil {
		return err
	}

	if err := sf.writer.WriteByte('\n'); err != nil {
		return err
	}

	err = sf.writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (sf *safeFile) Close() error {
	return sf.file.Close()
}

func (sf *safeFile) LoadAll() ([]models.ShortURL, error) {
	shortURLs := make([]models.ShortURL, 0)
	for sf.scanner.Scan() {
		data := sf.scanner.Bytes()
		var shortURL models.ShortURL
		err := json.Unmarshal(data, &shortURL)
		if err != nil {
			return nil, err
		}
		shortURLs = append(shortURLs, shortURL)
	}

	if err := sf.scanner.Err(); err != nil {
		return nil, err
	}

	return shortURLs, nil
}
