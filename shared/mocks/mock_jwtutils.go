package mocks

import "github.com/stretchr/testify/mock"

type MockJWTUtils struct {
	mock.Mock
}

func (m *MockJWTUtils) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
