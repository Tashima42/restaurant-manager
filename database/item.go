package database

import (
	"time"
)

type Item struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name" validate:"required"`
	Description string    `db:"description" json:"description" validate:"required"`
	Picture     string    `db:"picture" json:"picture" validate:"required"`
	Price       int       `db:"price" json:"price" validate:"required"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
