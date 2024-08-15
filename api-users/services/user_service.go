package services

import (
	"github.com/dieg0code/shared/json/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/dieg0code/shared/models"
)

type UserService interface {
	RegisterUser(createUserReq request.CreateUserRequest) (models.User, error)
	GetAllUsers() ([]response.UserResponse, error)
	GetUserByID(id string) (response.UserResponse, error)
}
