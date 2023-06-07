package controllers

import (
	"auction/database"
	"auction/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

type Product struct {
}

func (Product) ProductStore(c *fiber.Ctx) error {
	products := models.Product{}
	ProductFiles := models.Files{}
	products.StartDate = time.Now()
	products.EndDate = time.Now().AddDate(0, 0, 7)
	fileExept := []uint8{1} //ENUM - 1:image, 2:doc, 3:video, 4:zip
	if err := c.BodyParser(&products); err != nil {
		c.Status(401)
		return c.JSON(map[string]interface{}{
			"message": err,
		})
	}
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}
	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Image is required",
		})
	}
	if len(files) > 5 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Max image is 5",
		})
	}
	err = ProductFiles.SaveFile("products", fileExept, c, files)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	products.Files = ProductFiles.File
	err = database.Conn.DB.Create(&products).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    products,
	})

}

func (Product) ProductUpdate(c *fiber.Ctx) error {
	return nil
}

func (Product) ProductDelete(c *fiber.Ctx) error {
	productId := c.Params("id")
	product := models.Product{}
	database.Conn.DB.Clauses(clause.Returning{}).Delete(&product, productId)
	if product.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Product not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Product deleted successfully",
		"data":    product,
	})
}

func (Product) ProductGetAll(c *fiber.Ctx) error {
	slug := c.Params("name")
	if slug == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Category not found",
		})
	}
	products := []models.Product{}
	//filter by category slug

	database.Conn.DB.Preload("Brand").
		Preload("Category").
		Preload("Bids").
		Preload("Files").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("categories.slug = ?", slug).
		Find(&products)
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    products,
	})
}

func (Product) ProductGetID(c *fiber.Ctx) error {
	productId := c.Params("id")
	product := models.Product{}
	err := database.Conn.DB.Preload("Files").Preload("Brand").Preload("Seller.Identification").Preload("Bids.Seller.Identification").Preload("Category").First(&product, productId).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Product not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    product,
	})
}

func (Product) ProductGetAllBySellerID(c *fiber.Ctx) error {
	sellerId := c.Params("sellerid")
	seller := models.Seller{}
	sellererr := database.Conn.DB.Preload("Products.Category").Preload("Products.Brand").Preload("Products.Files").Preload("Products.Bids.Seller.Identification").First(&seller, sellerId).Error
	if sellererr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Seller not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    seller,
	})
}
