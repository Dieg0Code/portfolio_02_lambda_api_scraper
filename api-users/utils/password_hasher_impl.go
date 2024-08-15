package utils

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasherImpl struct{}

// HashPassword implements PasswordHasher.
func (p *PasswordHasherImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("[PasswordHasherImpl.HashPassword] error hashing password")
		return "", err
	}

	return string(hashedPassword), nil
}

func NewPasswordHasher() PasswordHasher {
	return &PasswordHasherImpl{}
}
