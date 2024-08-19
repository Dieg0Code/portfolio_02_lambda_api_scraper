package scraper

import (
	"errors"
	"fmt"
	"regexp"
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
func (s *ScraperImpl) CleanPrice(price string) ([]int, error) {
	// Captura numeros entre 1 y 3 digitos d{1,3} seguidos de 0 o más
	// grupos de 3 digitos (?:\.\d{3})* antesedidos por un punto
	// EJ: 123.456.789 sería \d{1,3} = 123 y (?:\.\d{3})* = .456.789
	re := regexp.MustCompile(`\d{1,3}(?:\.\d{3})*`)

	// Encontrar todas las coincidencias
	matches := re.FindAllString(price, -1)
	if len(matches) == 0 {
		logrus.Error("no prices found in string")
		return nil, errors.New("no prices found in string")
	}

	var prices []int
	for _, match := range matches {
		// Remover los puntos separadores de miles y convertir a entero
		cleaned := strings.ReplaceAll(match, ".", "")
		price, err := strconv.Atoi(cleaned)
		if err != nil {
			logrus.WithError(err).Error("error converting price to int")
			return nil, errors.New("error converting price to int")
		}
		prices = append(prices, price)
	}

	return prices, nil
}

// ScrapeData implements Scraper.
func (s *ScraperImpl) ScrapeData(protocol string, baseURL string, maxPage int, category string) ([]models.Product, error) {
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
		url := fmt.Sprintf("%s://%s/%s/page/%d/", protocol, baseURL, category, i)
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
