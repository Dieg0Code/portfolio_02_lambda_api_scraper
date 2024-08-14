package service

import (
	"github.com/dieg0code/scraper/src/repository"
	"github.com/dieg0code/scraper/src/scraper"
	"github.com/dieg0code/shared/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ScraperServiceImpl struct {
	Scraper           scraper.Scraper
	ScraperRepository repository.ScraperRepository
}

// GetProducts implements ScraperService.
func (s *ScraperServiceImpl) GetProducts() (bool, error) {
	const baseURL string = "cugat.cl/categoria-producto"

	err := s.ScraperRepository.DeleteAll()
	if err != nil {
		logrus.WithError(err).Error("[ProductServiceImpl.GetProducts] Error deleting all products")
		return false, err
	}

	logrus.Info("[ProductServiceImpl.UpdateData] Scraping data started")

	for _, categoryInfo := range scraper.Categories {
		products, err := s.Scraper.ScrapeData(baseURL, categoryInfo.MaxPage, categoryInfo.Category)
		if err != nil {
			logrus.WithError(err).Error("[ProductServiceImpl.UpdateData] Error scraping data")
			return false, err
		}
		for _, product := range products {
			productModel := models.Product{
				ProductID:       uuid.New().String(),
				Name:            product.Name,
				Category:        product.Category,
				OriginalPrice:   product.OriginalPrice,
				DiscountedPrice: product.DiscountedPrice,
			}
			_, err := s.ScraperRepository.Create(productModel)
			if err != nil {
				logrus.WithError(err).Error("[ProductServiceImpl.UpdateData] Error creating product")
				return false, err
			}
		}
	}

	logrus.Info("[ProductServiceImpl.UpdateData] Data scraped successfully")
	return true, nil
}

func NewScraperServiceImpl(scraper scraper.Scraper, scraperRepository repository.ScraperRepository) ScraperService {
	return &ScraperServiceImpl{
		Scraper:           scraper,
		ScraperRepository: scraperRepository,
	}
}
