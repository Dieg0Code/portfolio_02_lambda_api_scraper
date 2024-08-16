package controllers

import (
	"fmt"

	"github.com/dieg0code/api-users/services"
	"github.com/dieg0code/shared/json/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserControllerImpl struct {
	userService services.UserService
}

// LogInUser implements UserController.
func (u *UserControllerImpl) LogInUser(c *gin.Context) {
	loginRequest := request.LogInUserRequest{}

	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		logrus.WithError(err).Error("[UserControllerImpl.LogInUser] Error binding JSON")
		errorResponse := response.BaseResponse{
			Code:    400,
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		}

		c.JSON(400, errorResponse)
		return
	}

	loginResponse, err := u.userService.LogInUser(loginRequest)
	if err != nil {
		logrus.WithError(err).Error("[UserControllerImpl.LogInUser] Error logging in user")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "error",
			Message: "Error logging in user",
			Data:    nil,
		}

		c.JSON(500, errorResponse)
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", loginResponse.Token))
	webResponse := response.BaseResponse{
		Code:    200,
		Status:  "success",
		Message: "User logged in successfully",
		Data:    loginResponse,
	}

	c.JSON(200, webResponse)
}

// GetAllUsers implements UserController.
func (u *UserControllerImpl) GetAllUsers(c *gin.Context) {
	users, err := u.userService.GetAllUsers()
	if err != nil {
		logrus.WithError(err).Error("[UserControllerImpl.GetAllUsers] Error getting all users")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "error",
			Message: "Error getting all users",
			Data:    nil,
		}

		c.JSON(500, errorResponse)
		return
	}

	webResponse := response.BaseResponse{
		Code:    200,
		Status:  "success",
		Message: "Success getting all users",
		Data:    users,
	}

	c.JSON(200, webResponse)
}

// GetUserByID implements UserController.
func (u *UserControllerImpl) GetUserByID(c *gin.Context) {
	userID := c.Param("userID")

	user, err := u.userService.GetUserByID(userID)
	if err != nil {
		logrus.WithError(err).Error("[UserControllerImpl.GetUserByID] Error getting user by ID")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "error",
			Message: "Error getting user by ID",
			Data:    nil,
		}

		c.JSON(500, errorResponse)
		return
	}

	webResponse := response.BaseResponse{
		Code:    200,
		Status:  "success",
		Message: "User fetched successfully",
		Data:    user,
	}

	c.JSON(200, webResponse)
}

// RegisterUser implements UserController.
func (u *UserControllerImpl) RegisterUser(c *gin.Context) {
	registerUserRequest := request.CreateUserRequest{}

	err := c.ShouldBindJSON(&registerUserRequest)
	if err != nil {
		logrus.WithError(err).Error("[UserControllerImpl.RegisterUser] Error binding JSON")
		errorResponse := response.BaseResponse{
			Code:    400,
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		}

		c.JSON(400, errorResponse)
		return
	}

	user, err := u.userService.RegisterUser(registerUserRequest)
	if err != nil {
		logrus.WithError(err).Error("[UserControllerImpl.RegisterUser] Error registering user")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "error",
			Message: "Error registering user",
			Data:    nil,
		}

		c.JSON(500, errorResponse)
		return
	}

	webResponse := response.BaseResponse{
		Code:    201,
		Status:  "success",
		Message: "User registered successfully",
		Data:    fmt.Sprintf("User %s registered successfully ID: %s", user.Username, user.UserID),
	}

	c.JSON(201, webResponse)
}

func NewUserControllerImpl(userService services.UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}
