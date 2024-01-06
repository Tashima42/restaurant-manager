package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tashima42/restaurant-manager/database"
)

func (cr *Controller) CreateOrder(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	order := &database.Order{}
	cr.Logger.Info(requestID, ": unmarshal request body")
	if err := json.Unmarshal(c.Body(), order); err != nil {
		return err
	}

	if err := cr.Validate.Struct(order); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	cr.Logger.Info(requestID, ": creating order")
	if err := database.CreateOrder(c.Context(), cr.DB, order); err != nil {
		return err
	}
	cr.Logger.Info(requestID, ": order created")
	return c.JSON(map[string]interface{}{"success": true})
}
