package models

type Brand struct {
	ID      int       `json:"id"`
	Name    string    `json:"name" form:"name"`
	Product []Product `json:"product,omitempty"`
	File    []File    `gorm:"polymorphic:Table;polymorphicValue:brands" json:"file"`
}
