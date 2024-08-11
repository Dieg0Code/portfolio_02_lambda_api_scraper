package service

import (
	"testing"

	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/serverles-api-scraper/api/data/response"
	"github.com/dieg0code/serverles-api-scraper/api/models"
	"github.com/dieg0code/serverles-api-scraper/api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

type MockScraper struct {
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

func (m *MockScraper) ScrapeData(baseURL string, maxPage int, category string) ([]models.Product, error) {
	args := m.Called(baseURL, maxPage, category)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockScraper) CleanPrice(price string) (int, error) {
	args := m.Called(price)
	return args.Int(0), args.Error(1)
}

func TestPoductService_GetAll(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockScraper := new(MockScraper)
	productService := NewProductServiceImpl(mockRepo, mockScraper)

	expectedProducts := []response.ProductResponse{
		{
			ProductID:       "test-id",
			Name:            "Test Product",
			Category:        "Test Category",
			OriginalPrice:   100,
			DiscountedPrice: 90,
		},
		{
			ProductID:       "test-id-2",
			Name:            "Test Product 2",
			Category:        "Test Category 2",
			OriginalPrice:   200,
			DiscountedPrice: 190,
		},
	}

	mockRepo.On("GetAll").Return([]models.Product{
		{
			ProductID:       "test-id",
			Name:            "Test Product",
			Category:        "Test Category",
			OriginalPrice:   100,
			DiscountedPrice: 90,
		},
		{
			ProductID:       "test-id-2",
			Name:            "Test Product 2",
			Category:        "Test Category 2",
			OriginalPrice:   200,
			DiscountedPrice: 190,
		},
	}, nil)

	products, err := productService.GetAll()

	assert.NoError(t, err, "Expected no error, GetAll() returned an error")
	assert.Equal(t, expectedProducts, products, "Expected products to be equal to the expected products")

	mockRepo.AssertExpectations(t)

}

func TestPoductService_GetByID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockScraper := new(MockScraper)
	productService := NewProductServiceImpl(mockRepo, mockScraper)

	expectedProduct := response.ProductResponse{
		ProductID:       "test-id",
		Name:            "Test Product",
		Category:        "Test Category",
		OriginalPrice:   100,
		DiscountedPrice: 90,
	}

	mockRepo.On("GetByID", "test-id").Return(models.Product{
		ProductID:       "test-id",
		Name:            "Test Product",
		Category:        "Test Category",
		OriginalPrice:   100,
		DiscountedPrice: 90,
	}, nil)

	product, err := productService.GetByID("test-id")

	assert.NoError(t, err, "Expected no error, GetByID() returned an error")
	assert.Equal(t, expectedProduct, product, "Expected product to be equal to the expected product")

	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateData(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockScraper := new(MockScraper)
	productService := NewProductServiceImpl(mockRepo, mockScraper)

	mockRepo.On("DeleteAll").Return(nil)

	utils.Categories = []utils.CategoryInfo{
		{MaxPage: 10, Category: "bebidas-alcoholicas"},
		{MaxPage: 5, Category: "alimentos-basicos"},
	}

	mockScraper.On("ScrapeData", "cugat.cl/categoria-producto", 10, "bebidas-alcoholicas").Return([]models.Product{
		{
			Name:            "Product 1",
			Category:        "bebidas-alcoholicas",
			OriginalPrice:   100,
			DiscountedPrice: 90,
		},
	}, nil)

	mockScraper.On("ScrapeData", "cugat.cl/categoria-producto", 5, "alimentos-basicos").Return([]models.Product{
		{
			Name:            "Product 2",
			Category:        "alimentos-basicos",
			OriginalPrice:   200,
			DiscountedPrice: 180,
		},
	}, nil)

	mockRepo.On("Create", mock.MatchedBy(func(product models.Product) bool {
		return product.Name == "Product 1" && product.Category == "bebidas-alcoholicas" ||
			product.Name == "Product 2" && product.Category == "alimentos-basicos"
	})).Return(models.Product{}, nil)

	updateDataRequest := request.UpdateDataRequest{
		UpdateData: true,
	}

	success, err := productService.UpdateData(updateDataRequest)

	assert.NoError(t, err, "Expected no error, UpdateData() returned an error")
	assert.True(t, success, "Expected success to be true")

	mockRepo.AssertExpectations(t)
	mockScraper.AssertExpectations(t)
}
