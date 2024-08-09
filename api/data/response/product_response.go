package response

type ProductResponse struct {
	ProductID       string `json:"product_id"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	OriginalPrice   int    `json:"original_price"`
	DiscountedPrice int    `json:"discounted_price"`
}
