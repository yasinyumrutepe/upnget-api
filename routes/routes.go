package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	Home(api)
	Bid(api)
	Product(api)
	Category(api)
	Brand(api)
	Authentication(api)
	Seller(api)

}
