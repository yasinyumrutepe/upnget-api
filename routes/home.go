package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Home(api fiber.Router) {
	home := api.Group("/home")
	home.Get("/", controllers.Home{}.Home)
}
