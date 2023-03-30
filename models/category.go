package models

type Category struct {
	ID       uint      `json:"id"`
	ParentID uint      `json:"parentId"`
	Name     string    `json:"name"`
	Product  []Product `json:"product"`
}
