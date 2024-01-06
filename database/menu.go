package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Menu struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" validate:"required"`
	Items     []Item    `db:"items" json:"items,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func CreateMenu(ctx context.Context, db *sqlx.DB, m *Menu) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := CreateMenuTxx(tx, m); err != nil {
		return err
	}
	return tx.Commit()
}

func CreateMenuTxx(tx *sqlx.Tx, m *Menu) error {
	id := uuid.New().String()
	query := "INSERT INTO menus(id, name, created_at, updated_at) VALUES($1, $2, $3, $4);"
	_, err := tx.Exec(query, id, m.Name, time.Now(), time.Now())
	return err
}

func GetMenus(ctx context.Context, db *sqlx.DB) ([]Menu, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	menus, err := GetMenusTxx(tx)
	if err != nil {
		return nil, err
	}
	return menus, tx.Commit()
}

func GetMenusTxx(tx *sqlx.Tx) ([]Menu, error) {
	menus := []Menu{}
	query := "SELECT id, name, created_at, updated_at FROM menus;"
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		menu := new(Menu)
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.CreatedAt, &menu.UpdatedAt); err != nil {
			return nil, err
		}
		menus = append(menus, *menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menus, nil
}

// func GetMenuTxx(tx *sqlx.Tx) ([]Item, error) {
// 	items := []Item{}
// 	query := "SELECT id, name, description, picture, price, created_at, updated_at FROM items;"
// 	rows, err := tx.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for rows.Next() {
// 		item := new(Item)
// 		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Picture, &item.Price, &item.CreatedAt, &item.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, *item)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }
