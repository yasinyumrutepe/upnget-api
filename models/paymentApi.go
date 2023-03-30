package models

type PaymentApi struct {
	ID        int    `json:"id"`
	APIKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}
