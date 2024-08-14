package repository

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDynamoDB is a mock implementation of the DynamoDBAPI interface
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDB) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

// NewProductRepositoryImpl creates a new ProductRepositoryImpl with the given DynamoDBAPI and table name
func TestProductRepositoryImpl_GetByID(t *testing.T) {
	mockDB := new(MockDynamoDB)
	repo := NewProductRepositoryImpl(mockDB, "test-table")

	expectedProduct := models.Product{
		ProductID:       "test-id",
		Name:            "Test Product",
		Category:        "Test Category",
		OriginalPrice:   100,
		DiscountedPrice: 90,
	}

	mockDB.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"ProductID":       {S: &expectedProduct.ProductID},
			"Name":            {S: &expectedProduct.Name},
			"Category":        {S: &expectedProduct.Category},
			"OriginalPrice":   {N: stringPtr("100")},
			"DiscountedPrice": {N: stringPtr("90")},
		},
	}, nil)

	product, err := repo.GetByID("test-id")

	assert.NoError(t, err, "Expected no error, GetByID() returned an error")
	assert.Equal(t, expectedProduct, product, "Expected product to be equal to the expected product")
	mockDB.AssertExpectations(t)
}

// NewProductRepositoryImpl creates a new ProductRepositoryImpl with the given DynamoDBAPI and table name
func TestProductRepositoryImpl_GetAll(t *testing.T) {
	mockDB := new(MockDynamoDB)
	repo := NewProductRepositoryImpl(mockDB, "test-table")

	expectedProducts := []models.Product{
		{
			ProductID:       "test-id-1",
			Name:            "Test Product 1",
			Category:        "Test Category 1",
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
	}

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"ProductID":       {S: &expectedProducts[0].ProductID},
				"Name":            {S: &expectedProducts[0].Name},
				"Category":        {S: &expectedProducts[0].Category},
				"OriginalPrice":   {N: stringPtr("100")},
				"DiscountedPrice": {N: stringPtr("90")},
			},
			{
				"ProductID":       {S: &expectedProducts[1].ProductID},
				"Name":            {S: &expectedProducts[1].Name},
				"Category":        {S: &expectedProducts[1].Category},
				"OriginalPrice":   {N: stringPtr("200")},
				"DiscountedPrice": {N: stringPtr("180")},
			},
		},
	}, nil)

	products, err := repo.GetAll()

	assert.NoError(t, err, "Expected no error, GetAll() returned an error")
	assert.Equal(t, expectedProducts, products, "Expected products to be equal to the expected products")
	mockDB.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
