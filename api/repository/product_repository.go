package repository

import "github.com/dieg0code/serverles-api-scraper/api/models"

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id string) (models.Product, error)
	Create(product models.Product) (models.Product, error)
	DeleteAll() error
}
