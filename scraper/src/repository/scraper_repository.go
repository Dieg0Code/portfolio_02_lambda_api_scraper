package repository

import "github.com/dieg0code/shared/models"

type ScraperRepository interface {
	Create(product models.Product) (models.Product, error)
	DeleteAll() error
}
