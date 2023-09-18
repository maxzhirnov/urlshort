package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateUserID(t *testing.T) {
	auth := NewAuth()
	userId, err := auth.GenerateUserID(2)
	assert.NoError(t, err)
	assert.Greater(t, len(userId), 0)
}
