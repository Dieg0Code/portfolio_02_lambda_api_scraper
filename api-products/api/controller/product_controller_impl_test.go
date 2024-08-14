package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/serverles-api-scraper/api/data/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func TestProductController_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	productController := NewProductControllerImpl(mockService)

	router := gin.Default()
	router.GET("/products", productController.GetAll)

	mockService.On("GetAll").Return([]response.ProductResponse{
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
			DiscountedPrice: 180,
		},
	}, nil)

	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")

	var response response.BaseResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	assert.Equal(t, 200, response.Code, "Response code should be 200")
	assert.Equal(t, "OK", response.Status, "Response status should be Success")
}

func TestProductController_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	productController := NewProductControllerImpl(mockService)

	router := gin.Default()
	router.GET("/products/:productId", productController.GetByID)

	mockService.On("GetByID", "test-id").Return(response.ProductResponse{
		ProductID:       "test-id",
		Name:            "Test Product",
		Category:        "Test Category",
		OriginalPrice:   100,
		DiscountedPrice: 90,
	}, nil)

	req, err := http.NewRequest(http.MethodGet, "/products/test-id", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")

	var response response.BaseResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	assert.Equal(t, 200, response.Code, "Response code should be 200")
	assert.Equal(t, "OK", response.Status, "Response status should be Success")
}

func TestProductController_UpdateData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockProductService)
	productController := NewProductControllerImpl(mockService)

	router := gin.Default()
	router.PUT("/products", productController.UpdateData)

	mockService.On("UpdateData", request.UpdateDataRequest{
		UpdateData: true,
	}).Return(true, nil)

	reqBody, err := json.Marshal(request.UpdateDataRequest{
		UpdateData: true,
	})
	assert.NoError(t, err, "Expected no error marshalling request")

	req, err := http.NewRequest(http.MethodPut, "/products", bytes.NewBuffer(reqBody))
	assert.NoError(t, err, "Expected no error creating request")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")

	var response response.BaseResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshalling response")
	assert.Equal(t, 200, response.Code, "Response code should be 200")
	assert.Equal(t, "OK", response.Status, "Response status should be Success")
}
