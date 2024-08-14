package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dieg0code/shared/models"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type ScraperImpl struct {
	Collector *colly.Collector
}

// cleanPrice implements Scraper.
func (s *ScraperImpl) CleanPrice(price string) (int, error) {
	cleaned := strings.ReplaceAll(price, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	return strconv.Atoi(cleaned)
}

// scrapeData implements Scraper.
func (s *ScraperImpl) ScrapeData(baseURL string, maxPage int, category string) ([]models.Product, error) {
	var products []models.Product

	s.Collector.OnHTML(".product-small.box", func(e *colly.HTMLElement) {
		name := e.ChildText(".name.product-title a")
		category := e.ChildText(".category")
		originalPriceStr := e.ChildText(".price del .woocommerce-Price-amount.amount")
		discountPriceStr := e.ChildText(".price ins .woocommerce-Price-amount.amount")

		if discountPriceStr == "" {
			originalPriceStr = e.ChildText(".price .woocommerce-Price-amount.amount")
		}

		originalPrice, err := s.CleanPrice(originalPriceStr)
		if err != nil {
			originalPrice = 0
		}

		discountPrice, err := s.CleanPrice(discountPriceStr)
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
		url := fmt.Sprintf("https://%s/%s/page/%d/", baseURL, category, i)
		err := s.Collector.Visit(url)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to visit page %d at URL %s", i, url)
			if err.Error() == "Not Found" {
				continue
			}

			return nil, err
		}
	}

	return products, nil
}

func NewScraperImpl(collector *colly.Collector) *ScraperImpl {
	return &ScraperImpl{
		Collector: collector,
	}
}
