package database

import (
	"time"
)

type Cart struct {
	ID        string    `db:"id" json:"id"`
	Table     Table     `db:"table" json:"table" validate:"required"`
	Orders    []Order   `db:"orders" json:"orders"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
