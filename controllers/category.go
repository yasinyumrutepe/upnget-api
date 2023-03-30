package controllers

import (
	"auction/database"
	"auction/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

type Category struct{}

func (Category) CategoryStore(c *fiber.Ctx) error {
	category := models.Category{}
	if err := c.BodyParser(&category); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Category not created",
			"error":   "Invalid Data",
		})
	}
	err := database.Conn.DB.Create(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Category not created",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Category created",
		"data":    category,
	})
}
func (Category) CategoryUpdate(c *fiber.Ctx) error {
	category := models.Category{}
	if err := c.BodyParser(&category); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Category not updated",
			"error":   "Invalid Data",
		})
	}
	err := database.Conn.DB.First(&models.Category{
		ID: category.ID,
	}).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Category not updated",
			"error":   err.Error(),
		})
	}
	updateErr := database.Conn.DB.Model(&category).Updates(category).Error
	if updateErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Category not updated",
			"error":   updateErr.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Category updated",
		"data":    category,
	})
}
func (Category) CategoryDelete(c *fiber.Ctx) error {
	categoryId := c.Params("id")
	category := models.Category{}
	database.Conn.DB.Clauses(clause.Returning{}).Delete(&category, categoryId)
	if category.ID == 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Category not deleted",
			"error":   "Category not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Category deleted",
		"data":    category,
	})
}

func (Category) CategoryGetAll(c *fiber.Ctx) error {
	category := []models.Category{}
	database.Conn.DB.Preload("Product").Find(&category)
	return c.Status(200).JSON(fiber.Map{
		"message": "Category Get All",
		"data":    category,
	})
}

func (Category) GetCategorySubCategories(c *fiber.Ctx) error {
	categoryId := c.Params("categoryid")
	category := []models.Category{}
	database.Conn.DB.Where("parent_id=?", categoryId).Find(&category)

	return c.Status(200).JSON(fiber.Map{
		"message": "Category By SubCategories",
		"data":    category,
	})
}
