package models

type Order struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"productId"`
	UserID      uint    `json:"userId"`
	Quantity    float64 `json:"quantity"`
	OrderTypeID uint    `json:"orderType"`
}
