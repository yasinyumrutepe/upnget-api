package models

type OrderType struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Order []Order `json:"order"`
}
