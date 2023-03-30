package models

import "time"

type Identification struct {
	ID            int       `json:"id"`
	SellerID      uint      `json:"seller_id"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	IdentiyNumber string    `json:"identiyNumber"`
	Gender        uint      `json:"gender"`
	BirthDate     time.Time `json:"birthDate"`
	Email         string    `json:"email"`
}
