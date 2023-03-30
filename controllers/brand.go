package controllers

import (
	"auction/database"
	"auction/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

type Brand struct{}

func (Brand) BrandGetAll(c *fiber.Ctx) error {
	brand := []models.Brand{}
	database.Conn.DB.Find(&brand)
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    brand,
	})
}
func (Brand) BrandStore(c *fiber.Ctx) error {
	brand := models.Brand{}
	fileExept := []uint8{1} //ENUM - 1:image, 2:doc, 3:video, 4:zip
	if err := c.BodyParser(&brand); err != nil {
		c.Status(401)
		return c.JSON(map[string]interface{}{
			"error": "Invalid data",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}
	files := form.File["image"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Image is required",
		})
	}
	if len(files) > 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Max image is 1",
		})
	}
	BrandFiles := models.Files{}
	err = BrandFiles.SaveFile("products", fileExept, c, files)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	brand.File = BrandFiles.File

	createErr := database.Conn.DB.Create(&brand).Error
	if createErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": createErr.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Successfull added brand",
		"data":    brand,
	})
}

func (Brand) BrandUpdate(c *fiber.Ctx) error {
	brand := models.Brand{}
	if err := c.BodyParser(&brand); err != nil {
		c.Status(401)
		return c.JSON(map[string]interface{}{
			"error": "Invalid data",
		})
	}
	err := database.Conn.DB.First(&models.Brand{
		ID: brand.ID,
	}).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	updateErr := database.Conn.DB.Model(&brand).Updates(brand).Error
	if updateErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": updateErr.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Successfull updated brand",
		"data":    brand,
	})
}

func (Brand) BrandDelete(c *fiber.Ctx) error {
	brandID := c.Params("id")
	brand := models.Brand{}
	database.Conn.DB.Clauses(clause.Returning{}).Where("id=?", brandID).Delete(&brand)
	if brand.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Successfull deleted brand",
		"data":    brand,
	})
}

func (Brand) BrandGetProduct(c *fiber.Ctx) error {
	brand := models.Brand{}
	brandID := c.Params("id")
	err := database.Conn.DB.Preload("Products").First(&brand, brandID).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Brand Not Found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    brand,
	})
}
