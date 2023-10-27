package dtos

type Carts struct {
	ProductID  int64  `json:"product_id"`
	Quantity   int32  `json:"quantity"`
	SizeTypeID int8   `json:"size_type_id"`
	Comments   string `json:"comments"`
}
