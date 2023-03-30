package controllers

import (
	"auction/database"
	"auction/globals"
	"auction/models"

	"github.com/gofiber/fiber/v2"
)

type Bid struct {
}

func (Bid) Store(c *fiber.Ctx) error {
	bid := models.Bid{}
	productLastBid := models.Bid{}
	product := models.Product{}

	if err := c.BodyParser(&bid); err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error": "Invalid data",
		})
	}
	notVal := globals.ValidateStruct(&bid)
	if notVal != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": notVal,
		})
	}
	//Get Product
	err := database.Conn.DB.First(product, bid.ProductID).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.Conn.DB.Where("product_id=?", bid.ProductID).Order("price desc").First(&productLastBid).Error
	if err != nil {
		bidErr := database.Conn.DB.Create(&bid).Error
		if bidErr != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": bidErr.Error(),
			})
		}
	} else {
		if bid.Price <= productLastBid.Price {
			return c.Status(400).JSON(fiber.Map{
				"error": "Bid price must be higher than last bid",
			})
		}
		bidErr := database.Conn.DB.Create(&bid).Error
		if bidErr != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": bidErr.Error(),
			})

		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    bid,
	})

}
func (Bid) GetAllBidByProductID(c *fiber.Ctx) error {
	productID := c.Params("productid")
	product := models.Product{}
	err := database.Conn.DB.Where("id = ?", productID).Preload("Seller").Preload("Bids").Preload("Files").Preload("Category").Preload("Brand").First(&product).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    product,
	})

}
func (Bid) GetAllBidByUserID(c *fiber.Ctx) error {
	return c.SendString("Bid Store")

}
