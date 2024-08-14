package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/serverles-api-scraper/api/data/response"
	"github.com/dieg0code/serverles-api-scraper/api/repository"
	"github.com/sirupsen/logrus"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
}

// GetAll implements ProductService.
func (p *ProductServiceImpl) GetAll() ([]response.ProductResponse, error) {
	result, err := p.ProductRepository.GetAll()
	if err != nil {
		logrus.WithError(err).Error("[ProductServiceImpl.GetAll] Error getting all products")
		return nil, err
	}

	var products []response.ProductResponse
	for _, product := range result {
		productResponse := response.ProductResponse{
			ProductID:       product.ProductID,
			Name:            product.Name,
			Category:        product.Category,
			DiscountedPrice: product.DiscountedPrice,
			OriginalPrice:   product.OriginalPrice,
		}

		products = append(products, productResponse)
	}

	return products, nil
}

// GetByID implements ProductService.
func (p *ProductServiceImpl) GetByID(productID string) (response.ProductResponse, error) {
	result, err := p.ProductRepository.GetByID(productID)
	if err != nil {
		logrus.WithError(err).Error("[ProductServiceImpl.GetByID] Error getting product by ID")
		return response.ProductResponse{}, err
	}

	productResponse := response.ProductResponse{
		ProductID:       result.ProductID,
		Name:            result.Name,
		Category:        result.Category,
		OriginalPrice:   result.OriginalPrice,
		DiscountedPrice: result.DiscountedPrice,
	}

	return productResponse, nil
}

// UpdateData implements ProductService.
func (p *ProductServiceImpl) UpdateData(updateData request.UpdateDataRequest) (bool, error) {
	if !updateData.UpdateData {
		return false, nil
	}

	// Crear una nueva sesión de AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1"), // Reemplaza con tu región de AWS
	})
	if err != nil {
		logrus.WithError(err).Error("[ProductServiceImpl.UpdateData] Error creating AWS session")
		return false, err
	}

	// Crear un nuevo cliente Lambda
	svc := lambda.New(sess)

	// Preparar la entrada para invocar la función Lambda
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("scraper"),
		InvocationType: aws.String("Event"),
	}

	// Invocar la función Lambda
	_, err = svc.Invoke(input)
	if err != nil {
		logrus.WithError(err).Error("[ProductServiceImpl.UpdateData] Error invoking lambda function")
		return false, err
	}

	return true, nil
}

func NewProductServiceImpl(productRepository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepository,
	}
}
