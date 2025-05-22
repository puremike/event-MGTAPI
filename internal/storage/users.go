package storage

import (
	"database/sql"

	"golang.org/x/net/context"
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"_"`
}

func (u *UserStore) CreateUser(ctx context.Context, user *User) error {
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

func (u *UserStore) GetUserByID(ctx context.Context, userId int) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT id, name, email, password FROM users WHERE id = $1`

	user := &User{}
	if err := u.db.QueryRowContext(ctx, query, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT id, name, email, password FROM users WHERE email = $1`

	user := &User{}
	if err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
