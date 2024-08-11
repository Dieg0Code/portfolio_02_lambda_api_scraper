package utils

import "github.com/dieg0code/serverles-api-scraper/api/models"

type Scraper interface {
	ScrapeData(baseURL string, maxPage int, category string) ([]models.Product, error)
	CleanPrice(price string) (int, error)
}
