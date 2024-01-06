package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Table struct {
	ID          string    `db:"id" json:"id"`
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
		return errors.Wrap(err, tx.Rollback().Error())
	}
	return tx.Commit()
}

func CreateTableTxx(tx *sqlx.Tx, t *Table) error {
	id := uuid.New().String()
	query := "INSERT INTO tables(id, name, description, created_at, updated_at) VALUES($1, $2, $3, $4, $5);"
	_, err := tx.Exec(query, id, t.Name, t.Description, time.Now(), time.Now())
	return err
}

func GetTables(ctx context.Context, db *sqlx.DB) ([]Table, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	tables, err := GetTablesTxx(tx)
	if err != nil {
		return nil, errors.Wrap(err, tx.Rollback().Error())
	}
	return tables, tx.Commit()
}

func GetTablesTxx(tx *sqlx.Tx) ([]Table, error) {
	tables := []Table{}
	query := "SELECT id, name, description, created_at, updated_at FROM tables;"
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		table := new(Table)
		if err := rows.Scan(&table.ID, &table.Name, &table.Description, &table.CreatedAt, &table.UpdatedAt); err != nil {
			return nil, err
		}
		tables = append(tables, *table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}
