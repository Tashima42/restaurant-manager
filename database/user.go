package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRole int

const (
	UserRoleStudent     UserRole = 1
	UserRoleCoordinator UserRole = 2
	UserRoleSecretary   UserRole = 3
	UserRoleInstructor  UserRole = 4
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" validate:"required"`
	Email     string    `db:"email" json:"email" validate:"required"`
	Password  string    `db:"password" json:"password" validate:"required"`
	Role      UserRole  `db:"role"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func GetUserByEmailTxx(tx *sqlx.Tx, email string) (*User, error) {
	var u User
	query := "SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email=$1 LIMIT 1;"
	err := tx.Get(&u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByIDTxx(tx *sqlx.Tx, id int64) (*User, error) {
	var u User
	query := "SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id=$1 LIMIT 1;"
	err := tx.Get(&u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
