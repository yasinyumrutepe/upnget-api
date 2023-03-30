package models

type Address struct {
	ID         uint   `json:"id"`
	CityID     uint   `json:"city_id"`
	DistrictID uint   `json:"district_id"`
	Detail     string `json:"detail"`
	ZipCode    string `json:"zip_code"`
	SellerID   uint   `json:"seller_id"`
}
