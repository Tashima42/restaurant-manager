package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, tx.Rollback().Error())
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
		return nil, errors.Wrap(err, tx.Rollback().Error())
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

func GetMenuByID(ctx context.Context, db *sqlx.DB, id string) (*Menu, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	menus, err := GetMenuByIDTxx(tx, id)
	if err != nil {
		return nil, errors.Wrap(err, tx.Rollback().Error())
	}
	return menus, tx.Commit()
}

func GetMenuByIDTxx(tx *sqlx.Tx, id string) (*Menu, error) {
	menu := Menu{Items: []Item{}}
	query := `SELECT
							m.id as id,
							m.name as name,
							m.created_at as created_at,
							m.updated_at as updated_at,
							IFNULL(i.id, "") as item_id,
							IFNULL(i.name, "") as item_name,
							IFNULL(i.description, "") as item_description,
							IFNULL(i.picture, "") as item_picture,
							IFNULL(i.price, 0) as item_price
						FROM menus m
						INNER JOIN menu_items mi
						LEFT JOIN items i 
						ON mi.menu_id = m.id
						WHERE m.id = $1;`
	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := new(Item)
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.CreatedAt, &menu.UpdatedAt, &item.ID, &item.Name, &item.Description, &item.Picture, &item.Price); err != nil {
			return nil, err
		}
		if item.ID != "" {
			menu.Items = append(menu.Items, *item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &menu, nil
}

func VerifyMenuExistsTxx(tx *sqlx.Tx, id string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM menus WHERE id = $1);"
	if err := tx.Get(&exists, query, id); err != nil {
		return false, err
	}
	return exists, nil
}

func VerifyMenuItemsExistsTxx(tx *sqlx.Tx, menuID, itemID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM menu_items WHERE menu_id = $1 AND item_id = $2);"
	if err := tx.Get(&exists, query, menuID, itemID); err != nil {
		return false, err
	}
	return exists, nil
}

func CreateMenuItemTxx(tx *sqlx.Tx, menuID, itemID string) error {
	id := uuid.New().String()
	query := "INSERT INTO menu_items(id, menu_id, item_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5);"
	_, err := tx.Exec(query, id, menuID, itemID, time.Now(), time.Now())
	return err
}
