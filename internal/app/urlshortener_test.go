package app

import (
	"errors"
	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStorage struct {
	SaveFunc func(url models.URL) error
	GetFunc  func(id string) (models.URL, error)
}

func (ms *mockStorage) Save(url models.URL) error {
	return ms.SaveFunc(url)
}

func (ms *mockStorage) Get(id string) (models.URL, error) {
	return ms.GetFunc(id)
}

func Test_Create(t *testing.T) {
	type want struct {
		id  string
		err error
	}

	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "happy path",
			input: "google.com",
			want: want{
				id:  "12345678",
				err: nil,
			},
		},
		{
			name:  "empty url",
			input: "",
			want: want{
				id:  "",
				err: errors.New("originalURL shouldn't be empty string"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mockStorage{
				SaveFunc: func(url models.URL) error {
					return nil
				},
			}
			app := NewURLShortener(storage)
			actualID, actualErr := app.Create(tt.input)
			assert.Equal(t, len(tt.want.id), len(actualID))
			assert.Equal(t, tt.want.err, actualErr)
		})
	}
}

func Test_Get(t *testing.T) {
	type want struct {
		url models.URL
		err error
	}

	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "happy path",
			input: "12345678",
			want: want{
				url: models.URL{
					OriginalURL: "example.com",
					ID:          "12345678",
				},
				err: nil,
			},
		},
		{
			name:  "empty input case",
			input: "",
			want: want{
				url: models.URL{},
				err: errors.New("id shouldn't be empty string"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mockStorage{
				GetFunc: func(id string) (models.URL, error) {
					return models.URL{
						OriginalURL: "example.com",
						ID:          id,
					}, nil
				},
			}
			app := NewURLShortener(storage)
			actualURL, actualErr := app.Get(tt.input)
			assert.Equal(t, tt.want.url.OriginalURL, actualURL.OriginalURL)
			assert.Equal(t, tt.want.url.ID, actualURL.ID)
			assert.Equal(t, tt.want.err, actualErr)

		})
	}
}
