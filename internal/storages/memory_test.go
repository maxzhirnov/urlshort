package storages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxzhirnov/urlshort/internal/models"
)

func TestSafeMap(t *testing.T) {
	m := NewMemoryStorage()

	type want struct {
		id  string
		url string
	}

	tests := []struct {
		name     string
		inputURL models.ShortURL
		inputID  string
		want     want
	}{
		{
			name: "good case",
			inputURL: models.ShortURL{
				OriginalURL: "test.com",
				ID:          "1",
			},
			inputID: "1",
			want: want{
				id:  "1",
				url: "test.com",
			},
		},
		{
			name: "zero values",
			inputURL: models.ShortURL{
				OriginalURL: "",
				ID:          "",
			},
			inputID: "",
			want: want{
				id:  "",
				url: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.InsertURL(context.Background(), tt.inputURL)
			if err != nil {
				return
			}
			urlObjLoaded, ok := m.GetURLByID(context.Background(), tt.inputID)

			assert.Equal(t, tt.want.url, m.m[tt.inputID])
			assert.Equal(t, tt.want.url, urlObjLoaded.OriginalURL)
			assert.Equal(t, true, ok)
		})
	}
}
