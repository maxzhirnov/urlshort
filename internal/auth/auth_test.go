package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateUserID(t *testing.T) {
	auth := NewAuth()
	userID := auth.GenerateUUID()
	assert.Greater(t, len(userID), 0)
}
