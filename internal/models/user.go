package models

import (
	"database/sql"
	"errors"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

const userColumns = `id, email, password, name, role, created_at, updated_at`

func scanUser(row *sql.Row) (*User, error) {
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow(
		`SELECT `+userColumns+` FROM users WHERE email = ?`, email,
	)
	return scanUser(row)
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow(
		`SELECT `+userColumns+` FROM users WHERE id = ?`, id,
	)
	return scanUser(row)
}

func CreateUser(db *sql.DB, email, hashedPassword, name string) (*User, error) {
	result, err := db.Exec(
		`INSERT INTO users (email, password, name) VALUES (?, ?, ?)`,
		email, hashedPassword, name,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(db, int(id))
}
