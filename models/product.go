package models

import "time"

type Product struct {
	ID          uint      `json:"id" form:"id"`
	Name        string    `json:"name" form:"name"`
	Description string    `json:"description" form:"description"`
	Price       float64   `json:"price" form:"price"`
	SellerID    uint      `json:"seller_id" form:"seller_id"` // SellerID is the foreign key for Seller
	BrandID     uint      `json:"brand_id" form:"brand_id"`
	CategoryID  uint      `json:"category_id" form:"category_id"`
	StartDate   time.Time `json:"start_date" form:"start_date"`
	EndDate     time.Time `json:"end_date" form:"end_date"`
	Order       Order     `json:"-" `
	Files       []File    `gorm:"polymorphic:Table;polymorphicValue:products" json:"files,omitempty"`
	Bids        []Bid     `json:"bids,omitempty"`
	Category    Category  `json:"category"`
	Brand       Brand     `json:"brand"`
	Seller      Seller    `json:"seller"`
}
