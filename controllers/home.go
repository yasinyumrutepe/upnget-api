package controllers

import (
	"auction/database"
	"auction/models"

	"github.com/gofiber/fiber/v2"
)

type Home struct{}
type ReturnData struct {
	Products    []models.Product     `json:"products"`
	Brands      []models.Brand       `json:"brands"`
	Advertising []models.Advertising `json:"advertising"`
}
type Query struct {
	Search   string `json:"s"`
	Category uint   `json:"c"`
	Brand    string `json:"b"`
	Page     uint   `json:"p"`
	Key      string `json:"k"`
}

func (Home) Home(c *fiber.Ctx) error {
	returnData := ReturnData{}
	query := Query{}
	if err := c.QueryParser(&query); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	tx := database.Conn.DB
	tx.Preload("File").Find(&returnData.Brands)
	tx.Preload("File").Limit(5).Find(&returnData.Advertising)
	tx.Order("id desc").Preload("Category").Preload("Brand").Preload("Files").Limit(18).Find(&returnData.Products)

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    returnData,
	})
}
