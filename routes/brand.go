package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Brand(api fiber.Router) {
	api = api.Group("/brand")
	api.Get("/", controllers.Brand{}.BrandGetAll)
	api.Get("/:id", controllers.Brand{}.BrandGetProduct)
	api.Post("/", controllers.Brand{}.BrandStore)
	api.Delete("/:id<int>", controllers.Brand{}.BrandDelete)


}
