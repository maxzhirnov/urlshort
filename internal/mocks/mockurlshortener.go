package mocks

import "github.com/maxzhirnov/urlshort/internal/models"

type MockURLShortenerService struct {
	CreateFunc func(originalURL string) (id string, err error)
	GetFunc    func(id string) (url models.URL, err error)
}

func (m *MockURLShortenerService) Create(originalURL string) (string, error) {
	return m.CreateFunc(originalURL)
}

func (m *MockURLShortenerService) Get(id string) (models.URL, error) {
	return m.GetFunc(id)
}
