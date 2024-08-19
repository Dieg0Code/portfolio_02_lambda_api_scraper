package mocks

import (
	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetAll() ([]response.ProductResponse, error) {
	args := m.Called()
	return args.Get(0).([]response.ProductResponse), args.Error(1)
}
func (m *MockProductService) GetByID(productID string) (response.ProductResponse, error) {
	args := m.Called(productID)
	return args.Get(0).(response.ProductResponse), args.Error(1)
}
func (m *MockProductService) UpdateData(updateData request.UpdateDataRequest) (bool, error) {
	args := m.Called(updateData)
	return args.Bool(0), args.Error(1)
}
