package memorystorage

import (
	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSafeMap(t *testing.T) {
	m := NewSafeMap()

	type want struct {
		id  string
		url string
	}

	tests := []struct {
		name     string
		inputUrl models.URL
		inputID  string
		want     want
	}{
		{
			name: "good case",
			inputUrl: models.URL{
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
			inputUrl: models.URL{
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
			m.Store(tt.inputUrl)
			urlObjLoaded, ok := m.Load(tt.inputID)

			assert.Equal(t, tt.want.url, m.m[tt.inputID])
			assert.Equal(t, tt.want.url, urlObjLoaded.OriginalURL)
			assert.Equal(t, true, ok)
		})
	}
}
