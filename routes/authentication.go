package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Authentication(api fiber.Router) {
	api.Post("/login", controllers.Authentication{}.Login)

}
