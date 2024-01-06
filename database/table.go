package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Table struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name" validate:"required"`
	Description string    `db:"description" json:"description" validate:"required"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

func CreateTable(ctx context.Context, db *sqlx.DB, t *Table) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := CreateTableTxx(tx, t); err != nil {
		return err
	}
	return tx.Commit()
}

func CreateTableTxx(tx *sqlx.Tx, t *Table) error {
	query := "INSERT INTO tables(name, description, created_at, updated_at) VALUES($1, $2, $3, $4);"
	_, err := tx.Exec(query, t.Name, t.Description, time.Now(), time.Now())
	return err
}
