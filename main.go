package main

import (
	"auction/database"
	"auction/routes"
	"auction/secret"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	secret.LoadEnv("main", true)
	database.Conn.Connect()
	// database.Rcn.NewRedis()
	// database.Conn.DropSchema("public")
	// db.DBConn.Migrate()
	app := fiber.New(fiber.Config{})
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Static("/storage", "./storage")

	app.Get("/mig", func(c *fiber.Ctx) error {
		database.Conn.DropSchema("public")
		database.Conn.Migrate()
		database.Conn.Seed()
		return c.SendString("public dropped and created")
	})

	routes.Setup(app)
	app.Listen(":8000")
}
