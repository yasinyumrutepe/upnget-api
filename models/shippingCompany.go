package models

type ShippingCompany struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TrackingUrl  string  `json:"trackingUrl"`
	CargoPrice   float64 `json:"cargoPrice"`
	DeliveryDate string  `json:"deliveryDate"`
}
