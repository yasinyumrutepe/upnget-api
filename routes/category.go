package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Category(api fiber.Router) {
	api = api.Group("/category")
	api.Get("/", controllers.Category{}.CategoryGetAll)
	api.Get("/:categoryid<int>", controllers.Category{}.GetCategorySubCategories)
	api.Post("/", controllers.Category{}.CategoryStore)
	api.Put("/", controllers.Category{}.CategoryUpdate)
	api.Delete("/:id<int>", controllers.Category{}.CategoryDelete)
}
