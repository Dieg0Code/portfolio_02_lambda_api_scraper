package repository

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dieg0code/shared/mocks"
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepositoryImpl_GetByEmail(t *testing.T) {
	t.Run("GetByEmail_Success", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		expectedResult := models.User{
			UserID:   "test-id",
			Email:    "testing@test.com",
			Username: "test",
			Password: "test",
			Role:     "user",
		}

		item, err := dynamodbattribute.MarshalMap(expectedResult)
		assert.NoError(t, err, "Expected no error marshalling map")

		mockDB.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{item},
		}, nil)

		result, err := repo.GetByEmail("testing@test.com")
		assert.NoError(t, err, "Expected no error getting user")
		assert.Equal(t, expectedResult, result, "Expected user to be equal")

		mockDB.AssertExpectations(t)
	})

	t.Run("GetByEmail_UserNotFound", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		mockDB.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{},
		}, nil)

		_, err := repo.GetByEmail("nonexistent@test.com")

		assert.Error(t, err, "Expected error getting user")
		assert.Equal(t, "user not found", err.Error(), "Expected error message to be 'user not found'")

		mockDB.AssertExpectations(t)
	})

	t.Run("GetByEmail_QueryError", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		mockDB.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{}, errors.New("error getting user"))

		_, err := repo.GetByEmail("testing@test.com")

		assert.Error(t, err, "Expected error getting user")
		assert.Equal(t, "error getting user", err.Error(), "Expected error message to be 'error getting user'")

		mockDB.AssertExpectations(t)
	})
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		user := models.User{
			UserID:   "test-id",
			Username: "test",
			Email:    "testing@email.com",
			Password: "test",
			Role:     "user",
		}

		mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

		result, err := repo.Create(user)
		assert.NoError(t, err, "Expected no error creating user")

		assert.Equal(t, user, result, "Expected user to be equal")

		mockDB.AssertExpectations(t)
	})

	t.Run("Create_Error", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		user := models.User{
			UserID:   "test-id",
			Username: "test",
			Email:    "testing@test.com",
			Password: "test",
			Role:     "user",
		}

		mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, errors.New("error creating user"))

		_, err := repo.Create(user)
		assert.Error(t, err, "Expected error creating user")

		assert.Equal(t, "error creating user", err.Error(), "Expected error message to be 'error creating user'")

		mockDB.AssertExpectations(t)
	})
}

func TestUserRepositoryImpl_GetAll(t *testing.T) {
	t.Run("GetAll_Success", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		expectedResult := []models.User{
			{
				UserID:   "test-id",
				Email:    "testing@test.com",
				Username: "test",
				Password: "test",
				Role:     "user",
			},
			{
				UserID:   "test-id-2",
				Email:    "test2@test.com",
				Username: "test2",
				Password: "test2",
				Role:     "user",
			},
		}

		item1, err := dynamodbattribute.MarshalMap(expectedResult[0])
		assert.NoError(t, err, "Expected no error marshalling map")

		item2, err := dynamodbattribute.MarshalMap(expectedResult[1])
		assert.NoError(t, err, "Expected no error marshalling map")

		mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{item1, item2},
		}, nil).Once()

		result, err := repo.GetAll()
		assert.NoError(t, err, "Expected no error getting users")

		assert.Equal(t, expectedResult, result, "Expected users to be equal")

		mockDB.AssertExpectations(t)
	})

	t.Run("GetAll_Error", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{}, errors.New("error getting users"))

		_, err := repo.GetAll()
		assert.Error(t, err, "Expected error getting users")

		assert.Equal(t, "error getting users", err.Error(), "Expected error message to be 'error getting users'")

		mockDB.AssertExpectations(t)
	})
}

func TestUserRepositoryImpl_GetByID(t *testing.T) {
	t.Run("GetByID_Success", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		expectedResult := models.User{
			UserID:   "test-id",
			Email:    "test@test.com",
			Username: "test",
			Password: "test",
			Role:     "user",
		}

		item, err := dynamodbattribute.MarshalMap(expectedResult)
		assert.NoError(t, err, "Expected no error marshalling map")

		mockDB.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: item,
		}, nil).Once()

		result, err := repo.GetByID("test-id")
		assert.NoError(t, err, "Expected no error getting user")

		assert.Equal(t, expectedResult, result, "Expected user to be equal")

		mockDB.AssertExpectations(t)

	})

	t.Run("GetByID_UserNotFound", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		mockDB.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{}, nil)

		_, err := repo.GetByID("nonexistent-id")
		assert.Error(t, err, "Expected error getting user")

		assert.Equal(t, "user not found", err.Error(), "Expected error message to be 'user not found'")

		mockDB.AssertExpectations(t)
	})

	t.Run("GetByID_Error", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewUserRepositoryImpl(mockDB, "test-table")

		mockDB.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{}, errors.New("error getting user"))

		_, err := repo.GetByID("test-id")
		assert.Error(t, err, "Expected error getting user")

		assert.Equal(t, "error getting user", err.Error(), "Expected error message to be 'error getting user'")

		mockDB.AssertExpectations(t)
	})

}
