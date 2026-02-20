// handler/dealHandler.go
package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"

	"marketplace_entity/models"
	"marketplace_entity/responses"
	"marketplace_entity/services"
)

// ----------------- GET -----------------

func GetDealByID(c *fiber.Ctx) error {
	dealID, err := strconv.Atoi(c.Params("id"))
	if err != nil || dealID <= 0 {
		return responses.Error(c, 400, "invalid id")
	}

	deal, err := services.GetDealByID(dealID)
	if err != nil {
		if err == services.ErrDealNotFound {
			return responses.Error(c, 404, "deal not found")
		}
		return responses.Error(c, 500, err.Error())
	}

	return responses.Success(c, deal)
}

func GetDeals(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "100"))
	if err != nil || limit <= 0 {
		limit = 100
	}

	items, err := services.GetDeals(limit)
	if err != nil {
		return responses.Error(c, 500, err.Error())
	}

	return responses.Success(c, items)
}

// ----------------- POST -----------------

func CreateDeal(c *fiber.Ctx) error {
	var newDeal models.Deal

	if err := c.BodyParser(&newDeal); err != nil {
		return responses.Error(c, 400, "invalid json")
	}

	if newDeal.RequestID <= 0 || newDeal.ExecutorID <= 0 {
		return responses.Error(c, 400, "request_id and executor_id must be > 0")
	}

	if newDeal.Status == "" {
		newDeal.Status = "active"
	}

	createdDeal, err := services.CreateDeal(newDeal)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return responses.Error(c, 409, "request_id already exists")
			}
			if pgErr.Code == "23503" {
				return responses.Error(c, 400, "request_id or executor_id does not exist")
			}
		}
		return responses.Error(c, 500, "internal server error")
	}

	return responses.Success(c, createdDeal)
}

// ----------------- PUT -----------------

func UpdateDealFull(c *fiber.Ctx) error {
	dealID, err := strconv.Atoi(c.Params("id"))
	if err != nil || dealID <= 0 {
		return responses.Error(c, 400, "invalid id")
	}

	var updatedDeal models.Deal
	if err := c.BodyParser(&updatedDeal); err != nil {
		return responses.Error(c, 400, "invalid json")
	}

	resultDeal, err := services.UpdateDealFull(dealID, updatedDeal)
	if err != nil {
		if err == services.ErrDealNotFound {
			return responses.Error(c, 404, "deal not found")
		}
		return responses.Error(c, 500, err.Error())
	}

	return responses.Success(c, resultDeal)
}

// ----------------- PATCH -----------------

func UpdateDeal(c *fiber.Ctx) error {
	dealID, err := strconv.Atoi(c.Params("id"))
	if err != nil || dealID <= 0 {
		return responses.Error(c, 400, "invalid id")
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return responses.Error(c, 400, "invalid json")
	}

	updatedDeal, err := services.UpdateDeal(dealID, payload)
	if err != nil {
		if err == services.ErrDealNotFound {
			return responses.Error(c, 404, "deal not found")
		}
		return responses.Error(c, 400, err.Error())
	}

	return responses.Success(c, updatedDeal)
}

// ----------------- DELETE -----------------

func DeleteDeal(c *fiber.Ctx) error {
	dealID, err := strconv.Atoi(c.Params("id"))
	if err != nil || dealID <= 0 {
		return responses.Error(c, 400, "invalid id")
	}

	if err := services.DeleteDeal(dealID); err != nil {
		if err == services.ErrDealNotFound {
			return responses.Error(c, 404, "deal not found")
		}
		return responses.Error(c, 500, err.Error())
	}

	return responses.Success(c, fiber.Map{
		"message": "deal deleted",
		"id":      dealID,
	})
}
