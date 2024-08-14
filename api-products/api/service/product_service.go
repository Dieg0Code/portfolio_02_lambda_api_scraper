package service

import (
	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/serverles-api-scraper/api/data/response"
)

type ProductService interface {
	GetAll() ([]response.ProductResponse, error)
	GetByID(productID string) (response.ProductResponse, error)
	UpdateData(updateData request.UpdateDataRequest) (bool, error)
}
