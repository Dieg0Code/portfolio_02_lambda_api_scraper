package repository

import "github.com/dieg0code/shared/models"

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id string) (models.Product, error)
}
