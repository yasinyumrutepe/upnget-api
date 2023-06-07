package models

import (
	"gorm.io/gorm"
)

type Brand struct {
	ID      int       `json:"id"`
	Name    string    `json:"name" form:"name"`
	Product []Product `json:"product,omitempty"`
	File    []File    `gorm:"polymorphic:Table;polymorphicValue:brands" json:"file"`
}

func (Brand) Seed(db *gorm.DB) {

	brands := []Brand{
		{
			Name: "Apple",
		},
		{
			Name: "Samsung",
		},
		{
			Name: "Philips",
		},
		{
			Name: "Vestel",
		},
		{
			Name: "Dyson",
		},
		{
			Name: "Chopard",
		},
		{
			Name: "Mikimoto",
		},
		{
			Name: "Piaget",
		},
		{
			Name: "Antique",
		},
		{
			Name: "Borgonovo",
		},
		{
			Name: "Amefa",
		},
		{
			Name: "Art",
		},
		{
			Name: "Vintage",
		},
		{
			Name: "Other",
		},
	}

	db.Create(&brands)

	// if err != nil {
	// 	panic(err)
	// }
}
