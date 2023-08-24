package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxzhirnov/urlshort/internal/models"
)

func TestSafeMap(t *testing.T) {
	m := newSafeMap()

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
			m.Store(tt.inputURL)
			urlObjLoaded, ok := m.Load(tt.inputID)

			assert.Equal(t, tt.want.url, m.m[tt.inputID])
			assert.Equal(t, tt.want.url, urlObjLoaded.OriginalURL)
			assert.Equal(t, true, ok)
		})
	}
}