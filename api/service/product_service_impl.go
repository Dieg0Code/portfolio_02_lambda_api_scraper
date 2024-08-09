package service

import (
	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/serverles-api-scraper/api/data/response"
	"github.com/dieg0code/serverles-api-scraper/api/models"
	"github.com/dieg0code/serverles-api-scraper/api/repository"
	"github.com/dieg0code/serverles-api-scraper/api/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	Scraper           utils.Scraper
}

// GetAll implements ProductService.
func (p *ProductServiceImpl) GetAll() ([]response.ProductResponse, error) {
	result, err := p.ProductRepository.GetAll()
	if err != nil {
		logrus.WithError(err).Error("Error getting all products")
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
		logrus.WithError(err).Error("Error getting product by ID")
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
func (p *ProductServiceImpl) UpdateData(udateData request.UpdateDataRequest) (bool, error) {
	const BaseURL string = "cugat.cl/categoria-producto"

	if udateData.UpdateData {
		err := p.ProductRepository.DeleteAll()
		if err != nil {
			logrus.WithError(err).Error("Error deleting all products")
			return false, err
		}

		for _, categoryInfo := range utils.Categories {
			products, err := p.Scraper.ScrapeData(BaseURL, categoryInfo.MaxPage, categoryInfo.Category)
			if err != nil {
				logrus.WithError(err).Error("Error scraping data")
				return false, err
			}
			for _, product := range products {
				productModel := models.Product{
					ProductID:       uuid.New().String(),
					Name:            product.Name,
					Category:        product.Category,
					OriginalPrice:   product.OriginalPrice,
					DiscountedPrice: product.DiscountedPrice,
				}
				p.ProductRepository.Create(productModel)
			}
		}
	}

	return true, nil
}

func NewProductServiceImpl(productRepository repository.ProductRepository, scraper utils.Scraper) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepository,
		Scraper:           scraper,
	}
}
