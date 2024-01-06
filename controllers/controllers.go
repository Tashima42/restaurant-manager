package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Controller struct {
	DB        *sqlx.DB
	JWTSecret []byte
	Logger    *zap.SugaredLogger
	Validate  *validator.Validate
}
