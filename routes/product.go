package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Product(api fiber.Router) {
	api = api.Group("/product")
	api.Get("/", controllers.Product{}.ProductGetAll)
	api.Get("/:id", controllers.Product{}.ProductGetID)
	api.Post("/", controllers.Product{}.ProductStore)
	api.Put("/", controllers.Product{}.ProductUpdate)

}
