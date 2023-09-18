package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateUserID(t *testing.T) {
	auth := NewAuth()
	userId := auth.GenerateUUID()
	assert.Greater(t, len(userId), 0)
}
