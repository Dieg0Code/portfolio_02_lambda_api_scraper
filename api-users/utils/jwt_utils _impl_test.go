package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	t.Run("GenerateToken", func(t *testing.T) {
		// Create a new JWTUtils
		jwtUtils := NewJWTUtils()

		// Generate a token
		token, err := jwtUtils.GenerateToken("user_id")
		assert.NoError(t, err, "GenerateToken should not return an error")
		assert.NotEmpty(t, token, "Token should not be empty")
	})
}
