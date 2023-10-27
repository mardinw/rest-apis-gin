package deliveries

type Delivery struct {
	ID         int64  `json:"id"`
	RetailerID string `json:"retailer_id"`
	GroceryID  string `json:"grocery_id"`
	Status     string `json:"status"`
	Updated    int64  `json:"update"`
}
