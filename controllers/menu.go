package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

func (cr *Controller) AddItemToMenu(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	cr.Logger.Info(requestID, ": getting path params")
	menuID := c.Params("menu_id")
	if menuID == "" {
		return fiber.NewError(http.StatusBadRequest, "missing menu_id")
	}
	itemID := c.Params("item_id")
	if itemID == "" {
		return fiber.NewError(http.StatusBadRequest, "missing item_id")
	}

	tx, err := cr.DB.BeginTxx(c.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	menuExists, err := database.VerifyMenuExistsTxx(tx, menuID)
	if err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	if !menuExists {
		return fiber.NewError(http.StatusBadRequest, "item "+itemID+" doesn't exists")
	}
	itemExists, err := database.VerifyItemExistsTxx(tx, itemID)
	if err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	if !itemExists {
		return fiber.NewError(http.StatusBadRequest, "menu "+menuID+" doesn't exists")
	}
	menuItemsExists, err := database.VerifyMenuItemsExistsTxx(tx, menuID, itemID)
	if err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	if menuItemsExists {
		return fiber.NewError(http.StatusBadRequest, "item already is on the menu")
	}

	if err := database.CreateMenuItemTxx(tx, menuID, itemID); err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return c.JSON(map[string]interface{}{"success": true})
}

func (cr *Controller) GetMenu(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	menuID := c.Params("menu_id")
	if menuID == "" {
		return fiber.NewError(http.StatusBadRequest, "missing menu_id")
	}
	cr.Logger.Info(requestID, ": getting menu")
	menu, err := database.GetMenuByID(c.Context(), cr.DB, menuID)
	if err != nil {
		return err
	}
	return c.JSON(menu)
}
