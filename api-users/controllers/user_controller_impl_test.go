package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dieg0code/shared/json/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/dieg0code/shared/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogInUser(t *testing.T) {
	t.Run("LogInUser_Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.POST("/login", userController.LogInUser)

		userService.On("LogInUser", mock.Anything).Return(response.LogInUserResponse{
			Token: "token",
		}, nil)

		loginRequest := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "password",
		}

		reqBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err, "Expected no error marshalling request body")

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "success", response.Status, "Expected response status to be success")
		assert.Equal(t, "User logged in successfully", response.Message, "Expected response message to be 'User logged in successfully'")
		assert.Equal(t, "Bearer token", rec.Header().Get("Authorization"), "Expected Authorization header to be 'Bearer token'")
		assert.Equal(t, "token", response.Data.(map[string]interface{})["token"], "Expected response data token to be 'token'")
		userService.AssertExpectations(t)
	})

	t.Run("LogInUser_ErrorBindingJSON", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.POST("/login", userController.LogInUser)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("")))
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code 400")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "error", response.Status, "Expected response status to be error")
		assert.Equal(t, "Invalid request body", response.Message, "Expected response message to be 'Invalid request body'")

		userService.AssertExpectations(t)
	})

	t.Run("LogInUser_ErrorLoggingInUser", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.POST("/login", userController.LogInUser)

		userService.On("LogInUser", mock.Anything).Return(response.LogInUserResponse{}, assert.AnError)

		loginRequest := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "password",
		}

		reqBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err, "Expected no error marshalling request body")

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "error", response.Status, "Expected response status to be error")
		assert.Equal(t, "Error logging in user", response.Message, "Expected response message to be 'Error logging in user'")

		userService.AssertExpectations(t)
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("GetAllUsers_Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.GET("/users", userController.GetAllUsers)

		userService.On("GetAllUsers").Return([]response.UserResponse{
			{
				UserID:   "uuid",
				Email:    "test@test.com",
				Username: "test",
			},
			{
				UserID:   "uuid2",
				Email:    "test2@test.com",
				Username: "test2",
			},
		}, nil)

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "success", response.Status, "Expected response status to be success")
		assert.Equal(t, "Success getting all users", response.Message, "Expected response message to be 'Success getting all users'")
		assert.Equal(t, 2, len(response.Data.([]interface{})), "Expected response data to have 2 users")

		userService.AssertExpectations(t)

	})

	t.Run("GetAllUsers_ErrorGettingAllUsers", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.GET("/users", userController.GetAllUsers)

		userService.On("GetAllUsers").Return([]response.UserResponse{}, assert.AnError)

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected status code 500")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "error", response.Status, "Expected response status to be error")
		assert.Equal(t, "Error getting all users", response.Message, "Expected response message to be 'Error getting all users'")

		userService.AssertExpectations(t)
	})

}

func TestGetUserByID(t *testing.T) {
	t.Run("GetUserByID_Success", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.GET("/users/:userID", userController.GetUserByID)

		userID := "uuid"

		userService.On("GetUserByID", userID).Return(response.UserResponse{
			UserID:   "uuid",
			Username: "test",
			Email:    "test@test.com",
		}, nil)

		req, err := http.NewRequest(http.MethodGet, "/users/uuid", nil)
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Status code shoud be 200")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "success", response.Status, "Expected response status to be 'success'")

		userService.AssertExpectations(t)

	})

	t.Run("GetUserByID_Error", func(t *testing.T) {
		userService := new(mocks.MockUserService)
		userController := NewUserControllerImpl(userService)

		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.GET("/users/:userID", userController.GetUserByID)

		userID := "invalid-id"

		userService.On("GetUserByID", userID).Return(response.UserResponse{}, errors.New("error getting user"))

		req, err := http.NewRequest(http.MethodGet, "/users/invalid-id", nil)
		assert.NoError(t, err, "Expected no error creating request")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Status code shoud be 500")

		var response response.BaseResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Expected no error unmarshalling response body")
		assert.Equal(t, "error", response.Status, "Expected response status to be 'error'")

		userService.AssertExpectations(t)
	})
}
