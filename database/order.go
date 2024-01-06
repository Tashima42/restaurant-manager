package database

import (
	"time"
)

type Order struct {
	ID        int64     `db:"id" json:"id"`
	Item      Item      `db:"item" json:"item" validate:"required"`
	Quantity  int       `db:"quantity" json:"quantity" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
