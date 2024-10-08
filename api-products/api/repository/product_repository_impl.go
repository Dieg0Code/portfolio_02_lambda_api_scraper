package repository

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dieg0code/shared/models"
	"github.com/sirupsen/logrus"
)

type ProductRepositoryImpl struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
}

// GetAll implements ProductRepository.
func (p *ProductRepositoryImpl) GetAll() ([]models.Product, error) {
	input := &dynamodb.ScanInput{
		TableName: &p.tableName,
	}

	result, err := p.db.Scan(input)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.GetAll] error getting products")
		return nil, errors.New("error getting products")
	}

	var products []models.Product
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &products)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.GetAll] error unmarshalling products")
		return nil, errors.New("error getting products")
	}

	return products, nil
}

// GetByID implements ProductRepository.
func (p *ProductRepositoryImpl) GetByID(id string) (models.Product, error) {
	input := &dynamodb.GetItemInput{
		TableName: &p.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"ProductID": {
				S: aws.String(id),
			},
		},
	}

	result, err := p.db.GetItem(input)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.GetByID] error getting product")
		return models.Product{}, errors.New("error getting product")
	}

	if result.Item == nil {
		return models.Product{}, errors.New("product not found")
	}

	var product models.Product
	err = dynamodbattribute.UnmarshalMap(result.Item, &product)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.GetByID] error unmarshalling product")
		return models.Product{}, errors.New("error getting product")
	}

	return product, nil
}

func NewProductRepositoryImpl(db dynamodbiface.DynamoDBAPI, tableName string) ProductRepository {
	return &ProductRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}
