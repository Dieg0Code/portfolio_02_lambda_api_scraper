package repository

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dieg0code/shared/mocks"
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScraperRepository_Create(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewScraperRepositoryImpl(mockDB, "test-table")

		product := models.Product{
			ProductID:       "test-id",
			Name:            "Test Product",
			Category:        "Test Category",
			OriginalPrice:   100,
			DiscountedPrice: 90,
			LastUpdated:     "02-02-1996",
		}

		mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

		result, err := repo.Create(product)
		assert.NoError(t, err, "Expected no error creating product")

		assert.Equal(t, product, result, "Expected product to be the same")

		mockDB.AssertExpectations(t)
	})

	t.Run("Create_Error", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewScraperRepositoryImpl(mockDB, "test-table")

		product := models.Product{
			ProductID:       "test-id",
			Name:            "Test Product",
			Category:        "Test Category",
			OriginalPrice:   100,
			DiscountedPrice: 90,
			LastUpdated:     "02-02-1996",
		}

		mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, assert.AnError)

		result, err := repo.Create(product)
		assert.Error(t, err, "Expected error creating product")
		assert.Equal(t, models.Product{}, result, "Expected empty product")

		mockDB.AssertExpectations(t)
	})
}

func TestScraperRepository_DeleteAll(t *testing.T) {
	t.Run("DeleteAll_Success", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewScraperRepositoryImpl(mockDB, "test-table")

		// Mock scan result with items
		mockScanOutput := &dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{"ProductID": {S: aws.String("1")}},
				{"ProductID": {S: aws.String("2")}},
			},
		}

		mockDB.On("Scan", mock.Anything).Return(mockScanOutput, nil)
		mockDB.On("DeleteItem", mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)

		err := repo.DeleteAll()
		assert.NoError(t, err, "Expected no error deleting all products")

		// Ensure DeleteItem is called for each item
		mockDB.AssertNumberOfCalls(t, "DeleteItem", len(mockScanOutput.Items))
		mockDB.AssertExpectations(t)
	})

	t.Run("DeleteAll_NoProducts", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewScraperRepositoryImpl(mockDB, "test-table")

		// Mock scan result with no items
		mockScanOutput := &dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{},
		}

		mockDB.On("Scan", mock.Anything).Return(mockScanOutput, nil)

		err := repo.DeleteAll()
		assert.NoError(t, err, "Expected no error deleting all products")

		mockDB.AssertExpectations(t)
	})

	t.Run("DeleteAll_ErrorScanning", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewScraperRepositoryImpl(mockDB, "test-table")

		mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{}, assert.AnError)

		err := repo.DeleteAll()
		assert.Error(t, err, "Expected error deleting all products")

		mockDB.AssertExpectations(t)
	})

	t.Run("DeleteAll_ErrorDeleting", func(t *testing.T) {
		mockDB := new(mocks.MockDynamoDB)
		repo := NewScraperRepositoryImpl(mockDB, "test-table")

		// Mock scan result with items
		mockScanOutput := &dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{"ProductID": {S: aws.String("1")}},
				{"ProductID": {S: aws.String("2")}},
			},
		}

		mockDB.On("Scan", mock.Anything).Return(mockScanOutput, nil)
		mockDB.On("DeleteItem", mock.Anything).Return(&dynamodb.DeleteItemOutput{}, assert.AnError).Twice() // Se esperan dos llamadas

		err := repo.DeleteAll()
		assert.Error(t, err, "Expected error deleting all products")

		mockDB.AssertExpectations(t)
	})

}
