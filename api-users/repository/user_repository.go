package repository

import "github.com/dieg0code/shared/models"

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (models.User, error)
	Create(user models.User) (models.User, error)
}
