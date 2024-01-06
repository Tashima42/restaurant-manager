package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Item struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name" validate:"required_if=ID ''"`
	Description string     `db:"description" json:"description" validate:"required_if=ID ''"`
	Picture     string     `db:"picture" json:"picture" validate:"required_if=ID ''"`
	Price       int        `db:"price" json:"price" validate:"required_if=ID ''"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

func CreateItem(ctx context.Context, db *sqlx.DB, i *Item) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := CreateItemTxx(tx, i); err != nil {
		return errors.Wrap(err, tx.Rollback().Error())
	}
	return tx.Commit()
}

func CreateItemTxx(tx *sqlx.Tx, i *Item) error {
	id := uuid.New().String()
	query := "INSERT INTO items(id, name, description, picture, price, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7);"
	_, err := tx.Exec(query, id, i.Name, i.Description, i.Picture, i.Price, time.Now(), time.Now())
	return err
}

func GetItems(ctx context.Context, db *sqlx.DB) ([]Item, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	items, err := GetItemsTxx(tx)
	if err != nil {
		return nil, errors.Wrap(err, tx.Rollback().Error())
	}
	return items, tx.Commit()
}

func GetItemsTxx(tx *sqlx.Tx) ([]Item, error) {
	items := []Item{}
	query := "SELECT id, name, description, picture, price, created_at, updated_at FROM items;"
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := new(Item)
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Picture, &item.Price, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func GetItemByIDTxx(tx *sqlx.Tx, id string) (*Item, error) {
	var item Item
	query := "SELECT id, name, description, picture, price, created_at, updated_at FROM items WHERE id=$1 LIMIT 1;"
	err := tx.Get(&item, query, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func VerifyItemExistsTxx(tx *sqlx.Tx, id string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM items WHERE id = $1);"
	if err := tx.Get(&exists, query, id); err != nil {
		return false, err
	}
	return exists, nil
}
