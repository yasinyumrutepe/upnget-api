package models

import "time"

type Seller struct {
	ID             int            `json:"id"`
	UserID         uint           `json:"user_id"`
	User           User           `json:"user"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone"`
	DeletedAt      time.Time      `json:"deleted_at"`
	Identification Identification `json:"identification"`
	Address        []Address      `json:"address"`
	Products       []Product      `json:"product"`
}
