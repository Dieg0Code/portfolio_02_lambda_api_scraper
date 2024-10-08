package scraper

import "github.com/dieg0code/shared/models"

type Scraper interface {
	ScrapeData(protocol string, baseURL string, maxPage int, category string) ([]models.Product, error)
	CleanPrice(price string) ([]int, error)
}
