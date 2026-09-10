package model

type CartItem struct {
	BookId    uint  `json:"book_id"`
	Quantity  uint  `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
}
