package mocks

import (
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}
func (m *MockProductRepository) GetByID(id string) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(models.Product), args.Error(1)
}
func (m *MockProductRepository) Create(product models.Product) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}
func (m *MockProductRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}
