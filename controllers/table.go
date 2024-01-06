package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tashima42/restaurant-manager/database"
)

func (cr *Controller) CreateTable(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	table := &database.Table{}
	cr.Logger.Info(requestID, ": unmarshal request body")
	if err := json.Unmarshal(c.Body(), table); err != nil {
		return err
	}

	if err := cr.Validate.Struct(table); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	cr.Logger.Info(requestID, ": creating table")
	if err := database.CreateTable(c.Context(), cr.DB, table); err != nil {
		return err
	}
	cr.Logger.Info(requestID, ": table created")
	return c.JSON(map[string]interface{}{"success": true})
}
