package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
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

func CreateItem(ctx context.Context, db *sqlx.DB, i *Item) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := CreateItemTxx(tx, i); err != nil {
		return err
	}
	return tx.Commit()
}

func CreateItemTxx(tx *sqlx.Tx, i *Item) error {
	query := "INSERT INTO items(name, description, picture, price, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6);"
	_, err := tx.Exec(query, i.Name, i.Description, i.Picture, i.Price, time.Now(), time.Now())
	return err
}
