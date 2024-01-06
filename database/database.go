package database

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed schema_migrations
var schemaMigrations embed.FS

func Open(path string, migrateDownUp bool) (*sqlx.DB, error) {
	m, err := newMigrate(path)
	if err != nil {
		return nil, err
	}
	if migrateDownUp {
		if err := m.Down(); err != nil {
			return nil, err
		}
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return nil, err
		}
	}
	return sqlx.Open("sqlite", path)
}

func Close(db *sqlx.DB) error {
	return db.Close()
}

func newMigrate(path string) (*migrate.Migrate, error) {
	d, err := iofs.New(schemaMigrations, "schema_migrations")
	if err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("iofs", d, "sqlite://"+path)
}

func MigrateDown(ctx context.Context, path string) error {
	db, err := Open(path, true)
	if err != nil {
		return err
	}
	return Close(db)
}
