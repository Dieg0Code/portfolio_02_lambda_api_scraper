package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanPrice(t *testing.T) {
	scraper := NewScraperImpl(nil)

	// Test para un precio limpio
	priceStr := "$1.234"
	expectedPrice := 1234
	price, err := scraper.CleanPrice(priceStr)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, price)

	// Test para un string no numérico
	priceStr = "invalid"
	_, err = scraper.CleanPrice(priceStr)
	assert.Error(t, err)
}
