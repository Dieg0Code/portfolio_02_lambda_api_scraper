package mocks

import (
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/mock"
)

type MockScraper struct {
	mock.Mock
}

func (m *MockScraper) ScrapeData(protocol string, baseURL string, maxPage int, category string) ([]models.Product, error) {
	args := m.Called(protocol, baseURL, maxPage, category)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockScraper) CleanPrice(price string) ([]int, error) {
	args := m.Called(price)
	return args.Get(0).([]int), args.Error(1)
}
