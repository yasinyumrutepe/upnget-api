package models

type Bid struct {
	ID        uint    `json:"id"`
	Price     float64 `json:"price" validate:"required"`
	SellerID  uint    `json:"seller_id" validate:"required"`
	ProductID uint    `json:"product_id" validate:"required"`
	Product   Product `json:"product"`
	Seller    Seller  `json:"seller,omitempty"`
}
