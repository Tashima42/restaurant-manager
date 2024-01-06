package database

import (
	"time"
)

type Menu struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" validate:"required"`
	Items     []Item    `db:"items" json:"items"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
