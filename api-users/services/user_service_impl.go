package services

import (
	"errors"

	"github.com/dieg0code/api-users/repository"
	"github.com/dieg0code/api-users/utils"
	"github.com/dieg0code/shared/json/request"
	"github.com/dieg0code/shared/json/response"
	"github.com/dieg0code/shared/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validator      *validator.Validate
	PasswordHasher utils.PasswordHasher
	JWTUtils       utils.JWTUtils
}

// LogInUser implements UserService.
func (u *UserServiceImpl) LogInUser(logInUserReq request.LogInUserRequest) (response.LogInUserResponse, error) {
	err := u.Validator.Struct(logInUserReq)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.LogInUser] error validating login user request")
		return response.LogInUserResponse{}, err
	}

	user, err := u.UserRepository.GetByEmail(logInUserReq.Email)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.LogInUser] error getting user by email")
		return response.LogInUserResponse{}, err
	}

	err = u.PasswordHasher.ComparePassword(user.Password, logInUserReq.Password)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.LogInUser] error comparing password")
		return response.LogInUserResponse{}, errors.New("invalid password")
	}

	token, err := u.JWTUtils.GenerateToken(user.UserID)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.LogInUser] error generating token")
		return response.LogInUserResponse{}, err
	}

	logInUserResponse := response.LogInUserResponse{
		Token: token,
	}

	err = u.Validator.Struct(logInUserResponse)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.LogInUser] error validating login user response")
		return response.LogInUserResponse{}, err
	}

	return logInUserResponse, nil
}

// GetAllUsers implements UserService.
func (u *UserServiceImpl) GetAllUsers() ([]response.UserResponse, error) {
	users, err := u.UserRepository.GetAll()
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.GetAllUsers] error getting all users")
		return nil, err
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponse := response.UserResponse{
			UserID:   user.UserID,
			Username: user.Username,
			Email:    user.Email,
		}

		err := u.Validator.Struct(userResponse)
		if err != nil {
			logrus.WithError(err).Error("[UserServiceImpl.GetAllUsers] error validating user response")
			return nil, err
		}

		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}

// GetUserByID implements UserService.
func (u *UserServiceImpl) GetUserByID(id string) (response.UserResponse, error) {
	if id == "" {
		return response.UserResponse{}, errors.New("id is required")
	}

	user, err := u.UserRepository.GetByID(id)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.GetUserByID] error getting user by ID")
		return response.UserResponse{}, err
	}

	userResponse := response.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}

	err = u.Validator.Struct(userResponse)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.GetUserByID] error validating user response")
		return response.UserResponse{}, err
	}

	return userResponse, nil
}

// RegisterUser implements UserService.
func (u *UserServiceImpl) RegisterUser(createUserReq request.CreateUserRequest) (models.User, error) {
	err := u.Validator.Struct(createUserReq)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.RegisterUser] error validating create user request")
		return models.User{}, err
	}

	hashedPassword, err := u.PasswordHasher.HashPassword(createUserReq.Password)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.RegisterUser] error hashing password")
		return models.User{}, err
	}

	userModel := models.User{
		UserID:   uuid.New().String(),
		Username: createUserReq.Username,
		Email:    createUserReq.Email,
		Password: hashedPassword,
		Role:     createUserReq.Role,
	}

	user, err := u.UserRepository.Create(userModel)
	if err != nil {
		logrus.WithError(err).Error("[UserServiceImpl.RegisterUser] error creating user")
		return models.User{}, err
	}

	return user, nil
}

func NewUserServiceImpl(userRepository repository.UserRepository, validator *validator.Validate, passwordHaher utils.PasswordHasher, jwtUtils utils.JWTUtils) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validator:      validator,
		PasswordHasher: passwordHaher,
		JWTUtils:       jwtUtils,
	}
}
