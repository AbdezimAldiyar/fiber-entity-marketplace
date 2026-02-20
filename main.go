package main

import (
	"marketplace_entity/database"
	"marketplace_entity/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()

	routes.DealRoutes(app)
	app.Listen(":3000")
}
