package service

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/dieg0code/shared/mocks"
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPoductService_GetAll(t *testing.T) {
	t.Run("GetAll_Success", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

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
	})

	t.Run("GetAll_Error", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

		mockRepo.On("GetAll").Return([]models.Product{}, assert.AnError)

		products, err := productService.GetAll()

		assert.Error(t, err, "Expected error Getting all products")
		assert.Nil(t, products, "Expected products to be nil")
	})

}

func TestPoductService_GetByID(t *testing.T) {
	t.Run("GetByID_Success", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

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
	})

	t.Run("GetByID_Error", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

		mockRepo.On("GetByID", "test-id").Return(models.Product{}, assert.AnError)

		product, err := productService.GetByID("test-id")

		assert.Error(t, err, "Expected error Getting product by ID")
		assert.Equal(t, response.ProductResponse{}, product, "Expected product to be empty")
	})
}

func TestPoductService_UpdateData(t *testing.T) {
	t.Run("UpdateData_Success", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

		updateReq := request.UpdateDataRequest{
			UpdateData: true,
		}

		mockLambdaClient.On("Invoke", mock.Anything).Return(&lambda.InvokeOutput{}, nil)

		success, err := productService.UpdateData(updateReq)

		assert.NoError(t, err, "Expected no error, UpdateData() returned an error")
		assert.True(t, success, "Expected success to be true")

		mockLambdaClient.AssertExpectations(t)
	})

	t.Run("UpdateData_InvokeError", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

		updateReq := request.UpdateDataRequest{
			UpdateData: true,
		}

		mockLambdaClient.On("Invoke", mock.Anything).Return(&lambda.InvokeOutput{}, assert.AnError)

		success, err := productService.UpdateData(updateReq)

		assert.Error(t, err, "Expected error invoking lambda function")
		assert.False(t, success, "Expected success to be false")

		mockLambdaClient.AssertExpectations(t)
	})

	t.Run("UpdateData_NoUpdate", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		mockLambdaClient := new(mocks.MockLambdaClient)
		productService := NewProductServiceImpl(mockRepo, mockLambdaClient)

		updateReq := request.UpdateDataRequest{
			UpdateData: false,
		}

		success, err := productService.UpdateData(updateReq)

		assert.NoError(t, err, "Expected no error, UpdateData() returned an error")
		assert.False(t, success, "Expected success to be false")
	})
}
