package repository

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dieg0code/shared/models"
	"github.com/sirupsen/logrus"
)

type ScraperRepositoryImpl struct {
	db        dynamodbiface.DynamoDBAPI
	tableName string
}

// Create implements ScraperRepository.
func (s *ScraperRepositoryImpl) Create(product models.Product) (models.Product, error) {
	input := &dynamodb.PutItemInput{
		TableName: &s.tableName,
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

			"LastUpdated": {
				S: aws.String(product.LastUpdated),
			},
		},
	}

	_, err := s.db.PutItem(input)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.Create] error creating product")
		return models.Product{}, errors.New("error creating product")
	}

	return product, nil
}

// DeleteAll implements ScraperRepository.
func (s *ScraperRepositoryImpl) DeleteAll() error {
	scanInput := &dynamodb.ScanInput{
		TableName: &s.tableName,
	}

	result, err := s.db.Scan(scanInput)
	if err != nil {
		logrus.WithError(err).Error("[ProductRepositoryImpl.DeleteAll] error scanning products")
		return errors.New("error scanning products")
	}

	if len(result.Items) == 0 {
		logrus.Info("no products to delete")
		return nil
	}

	var deleteErrors []error

	// Iterar sobre los elementos escaneados y eliminarlos
	for _, item := range result.Items {
		deleteInput := &dynamodb.DeleteItemInput{
			TableName: &s.tableName,
			Key: map[string]*dynamodb.AttributeValue{
				"ProductID": item["ProductID"],
			},
		}

		_, err := s.db.DeleteItem(deleteInput)
		if err != nil {
			logrus.WithError(err).Error("[ProductRepositoryImpl.DeleteAll] error deleting product")
			deleteErrors = append(deleteErrors, err)
		}
	}

	// Si hubo errores, retornarlos combinados
	if len(deleteErrors) > 0 {
		logrus.WithError(err).Error("[ProductRepositoryImpl.DeleteAll] one or more errors occurred while deleting products")
		return errors.New("one or more errors occurred while deleting products")
	}

	return nil
}

func NewScraperRepositoryImpl(db dynamodbiface.DynamoDBAPI, tableName string) ScraperRepository {
	return &ScraperRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}
