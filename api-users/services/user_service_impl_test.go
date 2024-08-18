package services

import (
	"errors"
	"testing"

	"github.com/dieg0code/shared/json/request"
	"github.com/dieg0code/shared/mocks"
	"github.com/dieg0code/shared/models"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserServiceImpl_LogInUser(t *testing.T) {
	t.Run("LogInUser_Success", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		loginUserReq := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "password",
		}

		user := models.User{
			UserID:   "test-id",
			Email:    "test@test.com",
			Username: "test",
			Password: "password",
			Role:     "user",
		}

		userRepo.On("GetByEmail", loginUserReq.Email).Return(user, nil)
		passwordHasher.On("ComparePassword", user.Password, loginUserReq.Password).Return(nil)

		token := "test-token"

		jwtUtils.On("GenerateToken", user.UserID).Return(token, nil)

		logInUserResponse, err := userService.LogInUser(loginUserReq)
		assert.NoError(t, err, "Expected no error login user")

		assert.Equal(t, token, logInUserResponse.Token, "Expected token to be equal")

		userRepo.AssertExpectations(t)
	})

	t.Run("LogInUser_InvalidRequest", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		loginUserReq := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "",
		}

		logInUserResponse, err := userService.LogInUser(loginUserReq)

		assert.Error(t, err, "Expected error login user")
		assert.Empty(t, logInUserResponse, "Expected empty login user response")

		userRepo.AssertExpectations(t)
	})

	t.Run("LogInUser_GetByEmailError", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		loginUserReq := request.LogInUserRequest{
			Email:    "invalid-email@test.com",
			Password: "password",
		}

		userRepo.On("GetByEmail", loginUserReq.Email).Return(models.User{}, errors.New("error getting user by email"))

		logInUserResponse, err := userService.LogInUser(loginUserReq)

		assert.Error(t, err, "Expected error login user")
		assert.Empty(t, logInUserResponse, "Expected empty login user response")

		userRepo.AssertExpectations(t)
	})

	t.Run("LogInUser_ComparePasswordError", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		loginUserReq := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "wrong password",
		}

		user := models.User{
			UserID:   "test-id",
			Email:    "test@test.com",
			Username: "test",
			Password: "password",
			Role:     "user",
		}

		userRepo.On("GetByEmail", loginUserReq.Email).Return(user, nil)
		passwordHasher.On("ComparePassword", user.Password, loginUserReq.Password).Return(errors.New("error comparing password"))

		logInUserResponse, err := userService.LogInUser(loginUserReq)

		assert.Error(t, err, "Expected error login user")
		assert.Equal(t, "invalid password", err.Error(), "Expected invalid password error")
		assert.Empty(t, logInUserResponse, "Expected empty login user response")

		userRepo.AssertExpectations(t)
	})

	t.Run("LogInUser_GenerateTokenError", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		loginUserReq := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "password",
		}

		user := models.User{
			UserID:   "test-id",
			Email:    "test@test.com",
			Username: "test",
			Password: "password",
			Role:     "user",
		}

		userRepo.On("GetByEmail", loginUserReq.Email).Return(user, nil)
		passwordHasher.On("ComparePassword", user.Password, loginUserReq.Password).Return(nil)
		jwtUtils.On("GenerateToken", user.UserID).Return("", errors.New("error generating token"))

		logInUserResponse, err := userService.LogInUser(loginUserReq)

		assert.Error(t, err, "Expected error generating Token")
		assert.Empty(t, logInUserResponse, "Expected empty login user response")

		userRepo.AssertExpectations(t)

	})

	t.Run("LogInUser_InvalidResponse", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		loginUserReq := request.LogInUserRequest{
			Email:    "test@test.com",
			Password: "password",
		}

		user := models.User{
			UserID:   "test-id",
			Email:    "test@test.com",
			Username: "test",
			Password: "password",
			Role:     "user",
		}

		userRepo.On("GetByEmail", loginUserReq.Email).Return(user, nil)
		passwordHasher.On("ComparePassword", user.Password, loginUserReq.Password).Return(nil)
		jwtUtils.On("GenerateToken", user.UserID).Return("", nil)

		logInUserResponse, err := userService.LogInUser(loginUserReq)

		assert.Error(t, err, "Expected error validating login user response")
		assert.Empty(t, logInUserResponse, "Expected empty login user response")

		userRepo.AssertExpectations(t)

	})
}

func TestUserServiceImpl_GetAllUsers(t *testing.T) {

	t.Run("GetAllUser_Success", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		users := []models.User{
			{
				UserID:   "test-id",
				Email:    "test@test.com",
				Username: "test",
				Password: "test",
				Role:     "user",
			},
			{
				UserID:   "test-id-1",
				Email:    "test2@test.com",
				Username: "test2",
				Password: "test2",
				Role:     "user",
			},
		}

		userRepo.On("GetAll").Return(users, nil)

		userResponses, err := userService.GetAllUsers()

		assert.NoError(t, err, "Expected no error getting all users")
		assert.Len(t, userResponses, 2, "Expected 2 user responses")

		userRepo.AssertExpectations(t)
	})

	t.Run("GetAllUser_Error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()

		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		userRepo.On("GetAll").Return([]models.User{}, errors.New("error getting all users"))

		userResponses, err := userService.GetAllUsers()

		assert.Error(t, err, "Expected error getting all users")
		assert.Nil(t, userResponses, "Expected nil user responses")

		userRepo.AssertExpectations(t)
	})

	t.Run("GetAllUser_InvalidResponse", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()

		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		users := []models.User{
			{
				UserID:   "test-id",
				Email:    "invalid-email", // Invalid email format
				Username: "test",
				Password: "test",
				Role:     "user",
			},
			{
				UserID:   "test-id-1",
				Email:    "test2@test.com",
				Username: "", // Empty username
				Password: "test2",
				Role:     "user",
			},
		}

		userRepo.On("GetAll").Return(users, nil)

		userResponses, err := userService.GetAllUsers()

		assert.Error(t, err, "Expected error validating user response")
		assert.Nil(t, userResponses, "Expected nil user responses")
	})
}

func TestUserServiceImpl_GetUserByID(t *testing.T) {
	t.Run("GetUserByID_Success", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		userID := "test-id"

		user := models.User{
			UserID:   userID,
			Email:    "test@test.com",
			Username: "test",
			Password: "test",
			Role:     "user",
		}

		userRepo.On("GetByID", userID).Return(user, nil)

		userResponse, err := userService.GetUserByID(userID)

		assert.NoError(t, err, "Expected no error getting user by id")
		assert.Equal(t, user.UserID, userResponse.UserID, "Expected user id to be equal")

		userRepo.AssertExpectations(t)

	})

	t.Run("GetUserByID_Error", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		userID := "test-id"

		userRepo.On("GetByID", userID).Return(models.User{}, errors.New("error getting user by id"))

		userResponse, err := userService.GetUserByID(userID)

		assert.Error(t, err, "Expected error getting user by id")
		assert.Empty(t, userResponse, "Expected empty user response")

		userRepo.AssertExpectations(t)
	})

	t.Run("GetUserByID_InvalidResponse", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		userID := "test-id"

		user := models.User{
			UserID:   userID,
			Email:    "invalid-email", // Invalid email format
			Username: "test",
			Password: "test",
			Role:     "user",
		}

		userRepo.On("GetByID", userID).Return(user, nil)

		userResponse, err := userService.GetUserByID(userID)

		assert.Error(t, err, "Expected error validating user response")
		assert.Empty(t, userResponse, "Expected empty user response")

		userRepo.AssertExpectations(t)
	})
}

func TestUserServiceImpl_RegisterUser(t *testing.T) {
	t.Run("RegisterUser_Success", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		createUserReq := request.CreateUserRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "password",
			Role:     "user",
		}

		hashedPassword := "hashed-password"

		passwordHasher.On("HashPassword", createUserReq.Password).Return(hashedPassword, nil)

		userMatcher := mock.MatchedBy(func(user models.User) bool {
			return user.Username == createUserReq.Username &&
				user.Email == createUserReq.Email &&
				user.Password == hashedPassword &&
				user.Role == createUserReq.Role
		})

		expectedUser := models.User{
			UserID:   "test-id",
			Username: createUserReq.Username,
			Email:    createUserReq.Email,
			Password: hashedPassword,
			Role:     createUserReq.Role,
		}

		userRepo.On("Create", userMatcher).Return(expectedUser, nil)

		registeredUser, err := userService.RegisterUser(createUserReq)

		assert.NoError(t, err, "Expected no error registering user")
		assert.Equal(t, expectedUser.UserID, registeredUser.UserID, "Expected user id to be equal")

		userRepo.AssertExpectations(t)
	})

	t.Run("RegisterUser_InvalidRequest", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		createUserReq := request.CreateUserRequest{
			Username: "test",
			Email:    "invalid-email",
			Password: "password",
			Role:     "user",
		}

		registeredUser, err := userService.RegisterUser(createUserReq)

		assert.Error(t, err, "Expected error registering user")
		assert.Empty(t, registeredUser, "Expected empty registered user")

		userRepo.AssertExpectations(t)
	})

	t.Run("RegisterUser_HashPasswordError", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		createUserReq := request.CreateUserRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "password",
			Role:     "user",
		}

		passwordHasher.On("HashPassword", createUserReq.Password).Return("", errors.New("error hashing password"))

		registeredUser, err := userService.RegisterUser(createUserReq)

		assert.Error(t, err, "Expected error hashing password")
		assert.Empty(t, registeredUser, "Expected empty registered user")

		userRepo.AssertExpectations(t)
	})

	t.Run("RegisterUser_CreateUserError", func(t *testing.T) {
		userRepo := new(mocks.MockUserRepository)
		passwordHasher := new(mocks.MockPasswordHasher)
		jwtUtils := new(mocks.MockJWTUtils)
		validator := validator.New()
		userService := NewUserServiceImpl(userRepo, validator, passwordHasher, jwtUtils)

		createUserReq := request.CreateUserRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "password",
			Role:     "user",
		}

		hashedPassword := "hashed-password"

		passwordHasher.On("HashPassword", createUserReq.Password).Return(hashedPassword, nil)
		userRepo.On("Create", mock.Anything).Return(models.User{}, errors.New("error creating user"))

		registeredUser, err := userService.RegisterUser(createUserReq)

		assert.Error(t, err, "Expected error creating user")
		assert.Empty(t, registeredUser, "Expected empty registered user")

		userRepo.AssertExpectations(t)
	})
}
