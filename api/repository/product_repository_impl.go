package repository

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dieg0code/serverles-api-scraper/api/models"
	"github.com/sirupsen/logrus"
)

type ProductRepositoryImpl struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
}

// DeleteAll implements ProductRepository.
func (p *ProductRepositoryImpl) DeleteAll() error {
	// Escanear todos los elementos de la tabla
	scanInput := &dynamodb.ScanInput{
		TableName: &p.tableName,
	}

	result, err := p.db.Scan(scanInput)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.DeleteAll] error scanning products")
		return errors.New("error scanning products")
	}

	if len(result.Items) == 0 {
		logrus.Info("no products to delete")
		return nil
	}

	// Iterar sobre los elementos escaneados y eliminarlos
	for _, item := range result.Items {
		deleteInput := &dynamodb.DeleteItemInput{
			TableName: &p.tableName,
			Key: map[string]*dynamodb.AttributeValue{
				"ProductID": item["ProductID"],
			},
		}

		_, err := p.db.DeleteItem(deleteInput)
		if err != nil {
			logrus.WithError(err).Error("[ProductRepositoryImpl.DeleteAll] error deleting products")
			return errors.New("error deleting products")
		}
	}

	return nil
}

// Create implements ProductRepository.
func (p *ProductRepositoryImpl) Create(product models.Product) (models.Product, error) {
	input := &dynamodb.PutItemInput{
		TableName: &p.tableName,
		Item: map[string]*dynamodb.AttributeValue{
			"ProductID": {
				S: aws.String(product.ProductID),
			},
			"Name": {
				S: aws.String(product.Name),
			},
			"Category": {
				S: aws.String(product.Category),
			},
			"OriginalPrice": {
				N: aws.String(fmt.Sprintf("%d", product.OriginalPrice)),
			},
			"DiscountedPrice": {
				N: aws.String(fmt.Sprintf("%d", product.DiscountedPrice)),
			},
		},
	}

	_, err := p.db.PutItem(input)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.Create] error creating product")
		return models.Product{}, errors.New("error creating product")
	}

	return product, nil
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
