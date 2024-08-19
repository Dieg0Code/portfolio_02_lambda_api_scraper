package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/dieg0code/shared/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProductController_GetAll(t *testing.T) {
	t.Run("GetAll_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
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
	})

	t.Run("GetAll_Error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
		productController := NewProductControllerImpl(mockService)

		router := gin.Default()
		router.GET("/products", productController.GetAll)

		mockService.On("GetAll").Return([]response.ProductResponse{}, assert.AnError)

		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response")
		assert.Equal(t, 500, response.Code, "Response code should be 500")
		assert.Equal(t, "Internal Server Error", response.Status, "Response status should be Internal Server Error")

		mockService.AssertExpectations(t)
	})
}

func TestProductController_GetByID(t *testing.T) {
	t.Run("GetByID_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
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
	})

	t.Run("GetByID_Error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
		productController := NewProductControllerImpl(mockService)

		router := gin.Default()
		router.GET("/products/:productId", productController.GetByID)

		mockService.On("GetByID", "test-id").Return(response.ProductResponse{}, assert.AnError)

		req, err := http.NewRequest(http.MethodGet, "/products/test-id", nil)
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response")
		assert.Equal(t, 500, response.Code, "Response code should be 500")
		assert.Equal(t, "Internal Server Error", response.Status, "Response status should be Internal Server Error")

		mockService.AssertExpectations(t)
	})
}

func TestProductController_UpdateData(t *testing.T) {
	t.Run("UpdateData_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
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
	})

	t.Run("UpdateData_Failure", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
		productController := NewProductControllerImpl(mockService)

		router := gin.Default()
		router.PUT("/products", productController.UpdateData)

		mockService.On("UpdateData", request.UpdateDataRequest{
			UpdateData: false,
		}).Return(false, nil)

		reqBody, err := json.Marshal(request.UpdateDataRequest{
			UpdateData: false,
		})
		assert.NoError(t, err, "Expected no error marshalling request")

		req, err := http.NewRequest(http.MethodPut, "/products", bytes.NewBuffer(reqBody))
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 400")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response")
		assert.Equal(t, 400, response.Code, "Response code should be 400")
		assert.Equal(t, "Bad Request", response.Status, "Response status should be Bad Request")
		assert.Equal(t, "Error updating data", response.Message, "Response message should be UpdateData is required")

		mockService.AssertExpectations(t)
	})

	t.Run("UpdateData_ScrapingError", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockService := new(mocks.MockProductService)
		productController := NewProductControllerImpl(mockService)

		router := gin.Default()
		router.PUT("/products", productController.UpdateData)

		mockService.On("UpdateData", request.UpdateDataRequest{
			UpdateData: true,
		}).Return(false, assert.AnError)

		reqBody, err := json.Marshal(request.UpdateDataRequest{
			UpdateData: true,
		})
		assert.NoError(t, err, "Expected no error marshalling request")

		req, err := http.NewRequest(http.MethodPut, "/products", bytes.NewBuffer(reqBody))
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response")
		assert.Equal(t, 500, response.Code, "Response code should be 500")
		assert.Equal(t, "Internal Server Error", response.Status, "Response status should be Internal Server Error")

		mockService.AssertExpectations(t)
	})
}
