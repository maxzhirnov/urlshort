package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generateID(t *testing.T) {

	tests := []struct {
		name string
		val  int
		want int
	}{
		{
			name: "test 8 len",
			val:  8,
			want: 8,
		},
		{
			name: "test 35 len",
			val:  35,
			want: 35,
		},
		{
			name: "test 0 len",
			val:  0,
			want: 4,
		},
		{
			name: "test 0 len",
			val:  1,
			want: 4,
		},
		{
			name: "test 0 len",
			val:  2,
			want: 4,
		},
		{
			name: "test 0 len",
			val:  3,
			want: 4,
		},
		{
			name: "test 0 len",
			val:  4,
			want: 4,
		},
		{
			name: "test 0 len",
			val:  5,
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, len(generateID(tt.val)))
		})
	}
}

func Test_CheckURL(t *testing.T) {
	type want struct {
		IsValid bool
		URL     string
	}

	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "test with https",
			input: "https://ya.ru",
			want: want{
				IsValid: true,
				URL:     "https://ya.ru",
			},
		},
		{
			name:  "test with http",
			input: "http://ya.ru",
			want: want{
				IsValid: true,
				URL:     "http://ya.ru",
			},
		},
		{
			name:  "test without https and http",
			input: "ya.ru",
			want: want{
				IsValid: true,
				URL:     "http://ya.ru",
			},
		},
		{
			name:  "test not url",
			input: "yary",
			want: want{
				IsValid: false,
				URL:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, url := CheckURL(tt.input)
			assert.Equal(t, tt.want.URL, url)
			assert.Equal(t, tt.want.IsValid, isValid)
		})
	}
}
