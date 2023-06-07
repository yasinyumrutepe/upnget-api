package controllers

import (
	"auction/database"
	"auction/globals"
	"auction/models"

	"github.com/gofiber/fiber/v2"
)

type Seller struct {
}

func (Seller) GetSellerDetail(c *fiber.Ctx) error {
	seller_id := c.Params("sellerid")
	seller := models.Seller{}
	if seller_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Seller ID is required",
		})
	}
	err := database.Conn.DB.Preload("Identification").Preload("Address").First(&seller, seller_id).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Seller not found",
		})
	}
	return c.JSON(globals.Response{
		Error:   false,
		Message: "Seller found",
		Body:    seller,
	})

}
