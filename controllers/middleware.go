package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tashima42/restaurant-manager/database"
	"github.com/tashima42/restaurant-manager/hash"
)

type SignInUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GlobalError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (cr *Controller) SignIn(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	s := &SignInUser{}
	cr.Logger.Info(requestID, " unmarshal request body")
	if err := json.Unmarshal(c.Body(), s); err != nil {
		return err
	}

	cr.Logger.Info(requestID, " validate body")
	if err := cr.Validate.Struct(s); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	cr.Logger.Info(requestID, " starting transaction")
	tx, err := cr.DB.BeginTxx(c.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	cr.Logger.Info(requestID, " looking for user with email "+s.Email)
	user, err := database.GetUserByEmailTxx(tx, s.Email)
	if err != nil {
		cr.Logger.Info(requestID, " error: "+err.Error())
		if !strings.Contains(err.Error(), "no rows in result set") {
			return err
		}
		return fiber.NewError(http.StatusNotFound, "email "+s.Email+" not found")
	}

	cr.Logger.Info(requestID, " checking password")
	if !hash.CheckPassword(user.Password, s.Password) {
		return fiber.NewError(http.StatusUnauthorized, "incorrect password")
	}

	ac := hash.AuthClaims{}
	ac.User.ID = user.ID
	ac.User.Email = user.Email
	ac.User.Role = user.Role

	jwt, err := hash.NewJWT(cr.JWTSecret, ac)
	if err != nil {
		return errors.New("failed to generate jwt: " + err.Error())
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "auth-token"
	cookie.Value = jwt
	cookie.Expires = time.Now().Add(time.Hour * 24)
	c.Cookie(cookie)

	return c.JSON(map[string]interface{}{"token": jwt})
}

func (cr *Controller) ValidateToken(c *fiber.Ctx) error {
	requestID := fmt.Sprintf("%+v", c.Locals("requestid"))
	cr.Logger.Info(requestID, ": getting auth token cookie")
	token := c.Cookies("auth-token")
	if token == "" {
		authorizationHeader := c.GetReqHeaders()["authorization"]
		if len(authorizationHeader) <= 0 {
			return errors.New("missing authorization header value")
		}
		token = authorizationHeader[0]
	}
	if token == "" {
		cr.Logger.Info(requestID, ": missing auth token cookie")
		return fiber.NewError(http.StatusUnauthorized, "missing auth token")
	}
	cr.Logger.Info(requestID, ": parsing auth token")
	ac, err := hash.ParseJWT(cr.JWTSecret, token)
	if err != nil {
		cr.Logger.Error(requestID, err)
		return err
	}

	cr.Logger.Info(requestID, ": starting transaction")
	tx, err := cr.DB.BeginTxx(c.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	cr.Logger.Info(requestID, ": getting user")
	user, err := database.GetUserByIDTxx(tx, ac.User.ID)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	cr.Logger.Info(requestID, ": seting user on local variable")
	c.Locals("user", user)

	return c.Next()
}
func (cr *Controller) ErrorHandler(ctx *fiber.Ctx, err error) error {
	requestID := ctx.Locals("requestid")
	cr.Logger.Errorf("%s: %s", requestID, err.Error())
	code := fiber.StatusInternalServerError
	err = ctx.Status(code).JSON(GlobalError{Success: false, Message: err.Error()})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return nil
}
