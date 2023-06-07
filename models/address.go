package models

type Address struct {
	ID         uint   `json:"id"`
	CityID     string `json:"city"`
	DistrictID string `json:"district"`
	Detail     string `json:"detail"`
	ZipCode    string `json:"zip_code"`
	SellerID   uint   `json:"seller_id"`
}
