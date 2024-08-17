package models

type Product struct {
	ProductID       string `json:"product_id" dynamodbav:"ProductID"`
	Name            string `json:"name" dynamodbav:"Name"`
	Category        string `json:"category" dynamodbav:"Category"`
	OriginalPrice   int    `json:"original_price" dynamodbav:"OriginalPrice"`
	DiscountedPrice int    `json:"discounted_price" dynamodbav:"DiscountedPrice"`
	LastUpdated     string `json:"last_updated" dynamodbav:"LastUpdated"`
}
