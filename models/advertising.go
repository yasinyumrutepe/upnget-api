package models

type Advertising struct {
	ID   uint
	Link string `json:"link" form:"link"`
	Page string `json:"page" form:"page"`
	File []File `gorm:"polymorphic:Table;polymorphicValue:advertising" json:"file"`
}
