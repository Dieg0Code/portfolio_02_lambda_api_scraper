package service

import (
	"testing"

	"github.com/dieg0code/shared/mocks"
	"github.com/dieg0code/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScraperService_GetProducts(t *testing.T) {
	t.Run("GetProducts_Success", func(t *testing.T) {
		repo := new(mocks.MockScraperRepository)
		scraper := new(mocks.MockScraper)

		scraperService := NewScraperServiceImpl(scraper, repo)

		// Configurar los mocks
		repo.On("DeleteAll").Return(nil)
		scraper.On("ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything).Return([]models.Product{
			{
				Name:            "Product1",
				Category:        "Category1",
				OriginalPrice:   100,
				DiscountedPrice: 80,
			},
		}, nil)
		repo.On("Create", mock.Anything).Return(models.Product{}, nil)

		// Llamar a la función
		success, err := scraperService.GetProducts()

		// Verificar los resultados
		assert.NoError(t, err, "Expected no error, but got %v", err)
		assert.True(t, success, "Expected success to be true, but got %v", success)

		repo.AssertCalled(t, "DeleteAll")
		scraper.AssertCalled(t, "ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything)
		repo.AssertCalled(t, "Create", mock.Anything)
	})

	t.Run("GetProducts_ErrorDeletingAllProducts", func(t *testing.T) {
		repo := new(mocks.MockScraperRepository)
		scraper := new(mocks.MockScraper)

		scraperService := NewScraperServiceImpl(scraper, repo)

		// Configurar los mocks
		repo.On("DeleteAll").Return(assert.AnError)
		scraper.On("ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything).Return([]models.Product{
			{
				Name:            "Product1",
				Category:        "Category1",
				OriginalPrice:   100,
				DiscountedPrice: 80,
			},
		}, nil)
		repo.On("Create", mock.Anything).Return(models.Product{}, nil)

		// Llamar a la función
		success, err := scraperService.GetProducts()

		// Verificar los resultados
		assert.Error(t, err, "Expected an error, but got nil")
		assert.False(t, success, "Expected success to be false, but got %v", success)

		repo.AssertCalled(t, "DeleteAll")
		repo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("GetProducts_ErrorScrapingData", func(t *testing.T) {
		repo := new(mocks.MockScraperRepository)
		scraper := new(mocks.MockScraper)

		scraperService := NewScraperServiceImpl(scraper, repo)

		// Configurar los mocks
		repo.On("DeleteAll").Return(nil)
		scraper.On("ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything).Return([]models.Product{}, assert.AnError)
		repo.On("Create", mock.Anything).Return(models.Product{}, nil)

		// Llamar a la función
		success, err := scraperService.GetProducts()

		// Verificar los resultados
		assert.Error(t, err, "Expected an error, but got nil")
		assert.False(t, success, "Expected success to be false, but got %v", success)

		repo.AssertCalled(t, "DeleteAll")
		scraper.AssertCalled(t, "ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("GetProducts_ErrorCreatingProduct", func(t *testing.T) {
		repo := new(mocks.MockScraperRepository)
		scraper := new(mocks.MockScraper)

		scraperService := NewScraperServiceImpl(scraper, repo)

		// Configurar los mocks
		repo.On("DeleteAll").Return(nil)
		scraper.On("ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything).Return([]models.Product{
			{
				Name:            "Product1",
				Category:        "Category1",
				OriginalPrice:   100,
				DiscountedPrice: 80,
			},
		}, nil)
		repo.On("Create", mock.Anything).Return(models.Product{}, assert.AnError)

		// Llamar a la función
		success, err := scraperService.GetProducts()

		// Verificar los resultados
		assert.Error(t, err, "Expected an error, but got nil")
		assert.False(t, success, "Expected success to be false, but got %v", success)

		repo.AssertCalled(t, "DeleteAll")
		scraper.AssertCalled(t, "ScrapeData", "https", "cugat.cl/categoria-producto", mock.Anything, mock.Anything)
		repo.AssertCalled(t, "Create", mock.Anything)
	})
}
