package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Seller(api fiber.Router) {
	api = api.Group("/seller")
	api.Get("/:sellerid<int>", controllers.Seller{}.GetSellerDetail)

}
