package controllers

import (
	"auction/database"
	"auction/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Home struct{}
type ReturnData struct {
	ProductsFinishedToday []models.Product     `json:"products_finished_today"`
	ProductStartToday     []models.Product     `json:"products_start_today"`
	Brands                []models.Brand       `json:"brands"`
	Advertising           []models.Advertising `json:"advertising"`
}
type Query struct {
	Search string `json:"s"`
	Category uint `json:"c"`
	Brand string `json:"b"`
	Page uint `json:"p"`
	Key string `json:"k"`
}
func (Home) Home(c *fiber.Ctx) error {
	returnData := ReturnData{}
	query:=Query{}
	if err := c.QueryParser(&query); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	tx:=database.Conn.DB
	tx.Preload("File").Limit(5).Find(&returnData.Brands)
	tx.Preload("File").Limit(5).Find(&returnData.Advertising)
	
	currentTime := time.Now()
	lasthour := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 999, currentTime.Location())
	starthour := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	tx.Where("end_date<=?", lasthour).Preload("Bids",
		func(tx *gorm.DB) *gorm.DB {
			return tx.Order("price desc").Limit(1)

		}).Preload("Category").Preload("Brand").Preload("Files").Limit(5).Find(&returnData.ProductsFinishedToday)
		tx.Where("start_date>=?", starthour).Preload("Bids",
		func(tx *gorm.DB) *gorm.DB {
			return tx.Order("price desc").Limit(1)

		}).Preload("Category").Preload("Brand").Preload("Files").Limit(5).Find(&returnData.ProductStartToday)
	
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    returnData,
	})
}
