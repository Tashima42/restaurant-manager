package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Order struct {
	ID        string    `db:"id" json:"id"`
	Item      Item      `db:"item" json:"item" validate:"required"`
	Quantity  int       `db:"quantity" json:"quantity" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func CreateOrder(ctx context.Context, db *sqlx.DB, o *Order) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := CreateOrderTxx(tx, o); err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	return tx.Commit()
}

func CreateOrderTxx(tx *sqlx.Tx, o *Order) error {
	id := uuid.New().String()
	query := "INSERT INTO orders(id, item_id, quantity, created_at, updated_at) VALUES($1, $2, $3, $4, $5);"
	_, err := tx.Exec(query, id, o.Item.ID, o.Quantity, time.Now(), time.Now())
	return err
}
