package models

import (
	"gorm.io/gorm"
)

type Category struct {
	ID       uint      `json:"id"`
	ParentID uint      `json:"parentId"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Product  []Product `json:"product"`
}

func (Category) Seed(db *gorm.DB) {

	categories := []Category{
		{
			Name:     "Electronics",
			Slug:     "electronics",
			ParentID: 0,
		},
		{
			Name:     "Home & Living",
			Slug:     "home-living",
			ParentID: 0,
		},
		{
			Name:     "Art & Collectibles",
			Slug:     "art-collectibles",
			ParentID: 0,
		},
		{
			Name:     "Jewelry & Accessories",
			Slug:     "jewelry-accessories",
			ParentID: 0,
		},
		{
			Name:     "Antiques & Vintage",
			Slug:     "antiques-vintage",
			ParentID: 0,
		},
	}

	db.Create(&categories)

	// if err != nil {
	// 	panic(err)
	// }
}
