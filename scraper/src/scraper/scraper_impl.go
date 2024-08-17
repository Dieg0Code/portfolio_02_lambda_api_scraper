package scraper

import (
	"errors"
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

// CleanPrice implements Scraper.
// CleanPrice implements Scraper.
func (s *ScraperImpl) CleanPrice(price string) ([]int, error) {
	cleaned := strings.ReplaceAll(price, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.TrimSpace(cleaned) // Limpieza adicional

	if strings.Contains(cleaned, "\u2013") {
		priceParts := strings.Split(cleaned, "\u2013")
		var prices []int
		for _, part := range priceParts {
			part = strings.TrimSpace(part)
			if part == "" {
				logrus.WithError(errors.New("empty price part")).Error("empty price part")
				return nil, errors.New("empty price part")
			}
			price, err := strconv.Atoi(part)
			if err != nil {
				logrus.WithError(err).Error("error converting price to int")
				return nil, errors.New("error converting price to int")
			}

			prices = append(prices, price)
		}

		return prices, nil
	}

	priceInt, err := strconv.Atoi(cleaned)
	if err != nil {
		logrus.WithError(err).Error("error converting price to int")
		return nil, errors.New("error converting price to int")
	}

	return []int{priceInt}, nil
}

// ScrapeData implements Scraper.
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

		// Limpiar precios originales
		originalPrices, err := s.CleanPrice(originalPriceStr)
		if err != nil || len(originalPrices) == 0 {
			logrus.WithError(err).Error("error cleaning original price")
			originalPrices = []int{0}
		}

		// Limpiar precios con descuento
		discountPrices, err := s.CleanPrice(discountPriceStr)
		if err != nil || len(discountPrices) == 0 {
			logrus.WithError(err).Error("error cleaning discount price")
			discountPrices = []int{0}
		}

		// Crear una entrada por cada precio original
		for _, originalPrice := range originalPrices {
			for _, discountPrice := range discountPrices {
				product := models.Product{
					Name:            name,
					Category:        category,
					OriginalPrice:   originalPrice,
					DiscountedPrice: discountPrice,
				}
				products = append(products, product)
			}
		}
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

func NewScraperImpl(collector *colly.Collector) Scraper {
	return &ScraperImpl{
		Collector: collector,
	}
}
