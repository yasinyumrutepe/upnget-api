package routes

import (
	"auction/controllers"

	"github.com/gofiber/fiber/v2"
)

func Bid(api fiber.Router) {
	bid := api.Group("/bid")
	bid.Post("/", controllers.Bid{}.Store)
	bid.Get("/product/:productid<int>", controllers.Bid{}.GetAllBidByProductID)
	bid.Get("/seller/:sellerid<int>", controllers.Bid{}.GetAllBidByUserID)

}
