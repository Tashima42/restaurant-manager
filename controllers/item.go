package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tashima42/restaurant-manager/database"
)

func (cr *Controller) CreateItem(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	item := &database.Item{}
	cr.Logger.Info(requestID, ": unmarshal request body")
	if err := json.Unmarshal(c.Body(), item); err != nil {
		return err
	}

	if err := cr.Validate.Struct(item); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	cr.Logger.Info(requestID, ": creating item")
	if err := database.CreateItem(c.Context(), cr.DB, item); err != nil {
		return err
	}
	cr.Logger.Info(requestID, ": item created")
	return c.JSON(map[string]interface{}{"success": true})
}
