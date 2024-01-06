package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRole int

const (
	UserRoleManager UserRole = 1
	UserRoleWaiter  UserRole = 2
)

type User struct {
	ID        string    `db:"id" json:"id"`
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

func GetUserByIDTxx(tx *sqlx.Tx, id string) (*User, error) {
	var u User
	query := "SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id=$1 LIMIT 1;"
	err := tx.Get(&u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUserTxx(tx *sqlx.Tx, u *User) error {
	id := uuid.New().String()
	query := "INSERT INTO users(id, role, name, email, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7);"
	_, err := tx.Exec(query, id, u.Role, u.Name, u.Email, u.Password, time.Now(), time.Now())
	return err
}
