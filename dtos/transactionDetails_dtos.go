package dtos

type TransactionDetails struct {
	TransactionID     int64 `json:"transaction_id"`
	ProductID         int64 `json:"product_id"`
	Quantity          int8  `json:"quantity"`
	SizeTypeID        int8  `json:"size_type_id"`
	TotalPriceProduct int32 `json:"total_price_product"`
	OrderLineNumber   int8  `json:"order_line_num"`
}
