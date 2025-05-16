package storage

import "database/sql"

type UserModel struct {
	db *sql.DB
}

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"_"`
}
