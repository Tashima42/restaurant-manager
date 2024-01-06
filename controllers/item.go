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

func (cr *Controller) CreateItems(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	items := &[]database.Item{}
	cr.Logger.Info(requestID, ": unmarshal request body")
	if err := json.Unmarshal(c.Body(), items); err != nil {
		return err
	}

	for _, item := range *items {
		if err := cr.Validate.Struct(item); err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
	}

	cr.Logger.Info(requestID, " starting transaction")
	tx, err := cr.DB.BeginTxx(c.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	for _, item := range *items {
		cr.Logger.Info(requestID, ": creating item "+item.Name)
		if err := database.CreateItemTxx(tx, &item); err != nil {
			return errors.Wrap(err, tx.Rollback().Error())
		}
		cr.Logger.Info(requestID, ": created item "+item.Name)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return c.JSON(map[string]interface{}{"success": true})
}

func (cr *Controller) GetItems(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	cr.Logger.Info(requestID, ": getting items")
	items, err := database.GetItems(c.Context(), cr.DB)
	if err != nil {
		return err
	}
	return c.JSON(items)
}
