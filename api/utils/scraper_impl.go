package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dieg0code/serverles-api-scraper/api/models"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

type ScraperImpl struct{}

// cleanPrice implements Scraper.
func (s *ScraperImpl) cleanPrice(price string) (int, error) {
	cleaned := strings.ReplaceAll(price, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	return strconv.Atoi(cleaned)
}

// scrapeData implements Scraper.
func (s *ScraperImpl) ScrapeData(baseURL string, maxPage int, category string) ([]models.Product, error) {
	collector := colly.NewCollector()

	var products []models.Product

	collector.OnHTML(".product-small.box", func(e *colly.HTMLElement) {
		name := e.ChildText(".name.product-title a")
		category := e.ChildText(".category")
		originalPriceStr := e.ChildText(".price del .woocommerce-Price-amount.amount")
		discountPriceStr := e.ChildText(".price ins .woocommerce-Price-amount.amount")

		if discountPriceStr == "" {
			originalPriceStr = e.ChildText(".price .woocommerce-Price-amount.amount")
		}

		originalPrice, err := s.cleanPrice(originalPriceStr)
		if err != nil {
			originalPrice = 0
		}

		discountPrice, err := s.cleanPrice(discountPriceStr)
		if err != nil {
			discountPrice = 0
		}

		products = append(products, models.Product{
			Name:            name,
			Category:        category,
			OriginalPrice:   originalPrice,
			DiscountedPrice: discountPrice,
		})
	})

	for i := 1; i <= maxPage; i++ {
		logrus.Infof("Scraping page %d", i)
		collector.Visit(fmt.Sprintf("https://%s/%s/page/%d/", baseURL, category, i))
	}

	return products, nil
}

func NewScraperImpl() Scraper {
	return &ScraperImpl{}
}
