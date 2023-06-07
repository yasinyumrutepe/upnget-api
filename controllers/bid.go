package controllers

import (
	"auction/database"
	"auction/globals"
	"auction/models"

	"github.com/gofiber/fiber/v2"
)

type RequestBid struct {
	ProductID uint    `json:"product_id"`
	Price     float64 `json:"price"`
	SellerID  uint    `json:"seller_id"`
}
type Bid struct {
}

func (Bid) Store(c *fiber.Ctx) error {
	requestBid := RequestBid{}
	productLastBid := models.Bid{}
	product := models.Product{}

	if err := c.BodyParser(&requestBid); err != nil {
		c.Status(400)
		return c.JSON(map[string]interface{}{
			"error": "Invalid data",
		})
	}
	notVal := globals.ValidateStruct(&requestBid)
	if notVal != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": notVal,
		})
	}

	//Get Product
	err := database.Conn.DB.First(&product, requestBid.ProductID).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	addBid := models.Bid{
		ProductID: requestBid.ProductID,
		Price:     requestBid.Price,
		SellerID:  requestBid.SellerID,
	}
	err = database.Conn.DB.Where("product_id=?", requestBid.ProductID).Order("price desc").First(&productLastBid).Error
	if err != nil {
		bidErr := database.Conn.DB.Create(&addBid).Error
		if bidErr != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": bidErr.Error(),
			})
		}
	} else {
		if requestBid.Price <= productLastBid.Price {
			return c.Status(400).JSON(fiber.Map{
				"error": "Bid price must be higher than the last bid",
			})
		}
		bidErr := database.Conn.DB.Create(&addBid).Error
		if bidErr != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": bidErr.Error(),
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    addBid,
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
	seller_id := c.Params("sellerid")
	myBids := []models.Bid{}
	err := database.Conn.DB.Where("seller_id = ?", seller_id).Preload("Product").Preload("Product.Seller").Preload("Product.Bids").Preload("Product.Files").Preload("Product.Category").Preload("Product.Brand").Find(&myBids).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    myBids,
	})

}
