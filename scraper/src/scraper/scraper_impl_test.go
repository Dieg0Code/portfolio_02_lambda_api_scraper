package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dieg0code/shared/models"
	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

func TestCleanPrice(t *testing.T) {
	scraper := &ScraperImpl{}

	tests := []struct {
		input         string
		expected      []int
		expectedError bool
	}{
		{"123", []int{123}, false},
		{"1.234", []int{1234}, false},
		{"12.345.678", []int{12345678}, false},
		{"invalid price", nil, true},
		{"", nil, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			prices, err := scraper.CleanPrice(test.input)

			if test.expectedError {
				assert.Error(t, err, "Expected error cleaning price")
				assert.Nil(t, prices, "Expected nil prices")
			} else {
				assert.NoError(t, err, "Expected no error cleaning price")
				assert.Equal(t, test.expected, prices, "Expected cleaned prices to match")
			}
		})
	}
}

func TestScrapeData(t *testing.T) {
	t.Run("Scrape_Success", func(t *testing.T) {
		// Crear un servidor de prueba que devuelva un HTML estático
		ts := createTestServer()
		defer ts.Close()

		// Crear un nuevo scraper
		collector := colly.NewCollector()
		scraper := NewScraperImpl(collector)

		baseURL := strings.TrimPrefix(ts.URL, "http://")

		// Usa http:// para el servidor de prueba
		products, err := scraper.ScrapeData("http", baseURL, 1, "category")

		assert.NoError(t, err, "Expected no error scraping data")
		assert.Len(t, products, 1, "Expected 1 product")

		expectedProduct := models.Product{
			Name:            "Test Product",
			Category:        "category",
			OriginalPrice:   123456,
			DiscountedPrice: 7890,
		}

		assert.Equal(t, expectedProduct, products[0], "Expected product to match")
	})
}

// TestServer to simualte a real page to scrape
func createTestServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<div class="product-small box">
				<div class="name product-title"><a href="#">Test Product</a></div>
				<div class="category">category</div>
				<div class="price">
					<del><span class="woocommerce-Price-amount amount">123.456</span></del>
					<ins><span class="woocommerce-Price-amount amount">7.890</span></ins>
				</div>
			</div>
		`))
	})

	ts := httptest.NewServer(handler)
	return ts
}
