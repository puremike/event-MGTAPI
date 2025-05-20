package storage

import (
	"database/sql"

	"golang.org/x/net/context"
)

type UserModel struct {
	db *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"_"`
}

func (u *UserModel) CreateUser(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = tx.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (u *UserModel) GetUserByID(ctx context.Context, userId int) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT id, name, email FROM users WHERE id = $1`

	user := &User{}
	if err := u.db.QueryRowContext(ctx, query, userId).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
