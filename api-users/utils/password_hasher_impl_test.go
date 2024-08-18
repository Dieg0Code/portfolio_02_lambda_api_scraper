package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparePassword(t *testing.T) {
	t.Run("ComparePassword", func(t *testing.T) {
		// Create a new PasswordHasher
		passwordHasher := NewPasswordHasher()

		// Hash a password
		hashedPassword, err := passwordHasher.HashPassword("password")
		assert.NoError(t, err, "HashPassword should not return an error")

		// Compare the password
		err = passwordHasher.ComparePassword(hashedPassword, "password")
		assert.NoError(t, err, "ComparePassword should not return an error")

	})

	t.Run("ComparePasswordFail", func(t *testing.T) {
		// Create a new PasswordHasher
		passwordHasher := NewPasswordHasher()

		// Hash a password
		hashedPassword, err := passwordHasher.HashPassword("password")
		assert.NoError(t, err, "HashPassword should not return an error")

		// Compare the password
		err = passwordHasher.ComparePassword(hashedPassword, "wrong_password")
		assert.Error(t, err, "ComparePassword should return an error")
	})

}

func TestHashPassword(t *testing.T) {
	t.Run("HashPassword", func(t *testing.T) {
		// Create a new PasswordHasher
		passwordHasher := NewPasswordHasher()

		// Hash a password
		hashedPassword, err := passwordHasher.HashPassword("password")
		assert.NoError(t, err, "HashPassword should not return an error")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	})
}
