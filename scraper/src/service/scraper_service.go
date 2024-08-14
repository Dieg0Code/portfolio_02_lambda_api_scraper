package service

type ScraperService interface {
	GetProducts() (bool, error)
}
