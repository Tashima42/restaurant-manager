package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tashima42/restaurant-manager/database"
	"github.com/tashima42/restaurant-manager/hash"
	"go.uber.org/zap"
)

func (cr *Controller) CreateUser(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	user := &database.User{}
	cr.Logger.Info(requestID, " unmarshal request body")
	if err := json.Unmarshal(c.Body(), user); err != nil {
		return err
	}

	if err := cr.Validate.Struct(user); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	cr.Logger.Info(requestID, " starting transaction")
	tx, err := cr.DB.BeginTxx(c.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	cr.Logger.Info(requestID, " looking for user with email "+user.Email)
	if _, err := database.GetUserByEmailTxx(tx, user.Email); err != nil {
		cr.Logger.Info(requestID, " error: "+err.Error())
		if !strings.Contains(err.Error(), "no rows in result set") {
			return err
		}
		cr.Logger.Info(requestID, " user doesn't exists, continue")
	} else {
		zap.Error(errors.New(requestID + " email was already registered"))
		return fiber.NewError(http.StatusConflict, "email "+user.Email+" already was registered")
	}

	cr.Logger.Info(requestID, " hashing password")
	hashedPassword, err := hash.Password(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	cr.Logger.Info(requestID, " creating user")
	if err := database.CreateUserTxx(tx, user); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	cr.Logger.Info(requestID, " user created")
	return c.JSON(map[string]interface{}{"success": true})
}
