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
