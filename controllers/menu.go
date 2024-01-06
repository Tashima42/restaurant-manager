package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tashima42/restaurant-manager/database"
)

func (cr *Controller) CreateMenu(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	menu := &database.Menu{}
	cr.Logger.Info(requestID, ": unmarshal request body")
	if err := json.Unmarshal(c.Body(), menu); err != nil {
		return err
	}

	if err := cr.Validate.Struct(menu); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	cr.Logger.Info(requestID, ": creating menu")
	if err := database.CreateMenu(c.Context(), cr.DB, menu); err != nil {
		return err
	}
	cr.Logger.Info(requestID, ": menu created")
	return c.JSON(map[string]interface{}{"success": true})
}

func (cr *Controller) GetMenus(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	cr.Logger.Info(requestID, ": getting menus")
	menus, err := database.GetMenus(c.Context(), cr.DB)
	if err != nil {
		return err
	}
	return c.JSON(menus)
}
