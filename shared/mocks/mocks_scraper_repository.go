package mocks

import (
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/mock"
)

type MockScraperRepository struct {
	mock.Mock
}

func (m *MockScraperRepository) Create(product models.Product) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}
func (m *MockScraperRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}
