package routes

import (
	handler "marketplace_entity/handlers"

	"github.com/gofiber/fiber/v2"
)

func DealRoutes(app *fiber.App) {
	deals := app.Group("/deals")

	deals.Get("/", handler.GetDeals)
	deals.Get("/:id", handler.GetDealByID)

	deals.Post("/", handler.CreateDeal)
	deals.Put("/:id", handler.UpdateDealFull)
	deals.Patch("/:id", handler.UpdateDeal)
	deals.Delete("/:id", handler.DeleteDeal)

}
