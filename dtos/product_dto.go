package dtos

type Product struct {
	ProductCode string `json:"product_code,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Picture     string `json:"picture,omitempty"`
	Position    string `json:"position"`
	Quantity    int16  `json:"quantity,omitempty"`
	SizeTypeId  int    `json:"size_type_id,omitempty"`
	CategoryId  int    `json:"category_id,omitempty"`
	BuyPrice    int32  `json:"buy_price,omitempty"`
	MRP         int32  `json:"min_retail_price,omitempty"`
	Defective   int32  `json:"defective,omitempty"`
	Active      bool   `json:"active"`
}
