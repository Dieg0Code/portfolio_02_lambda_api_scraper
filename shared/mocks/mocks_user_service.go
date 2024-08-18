package mocks

import (
	"github.com/dieg0code/shared/json/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(createUserReq request.CreateUserRequest) (models.User, error) {
	args := m.Called(createUserReq)
	return args.Get(0).(models.User), args.Error(1)
}
func (m *MockUserService) GetAllUsers() ([]response.UserResponse, error) {
	args := m.Called()
	return args.Get(0).([]response.UserResponse), args.Error(1)
}
func (m *MockUserService) GetUserByID(id string) (response.UserResponse, error) {
	args := m.Called(id)
	return args.Get(0).(response.UserResponse), args.Error(1)
}
func (m *MockUserService) LogInUser(logInUserReq request.LogInUserRequest) (response.LogInUserResponse, error) {
	args := m.Called(logInUserReq)
	return args.Get(0).(response.LogInUserResponse), args.Error(1)
}
