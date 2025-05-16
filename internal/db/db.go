package db

import (
	"context"
	"database/sql"
	"time"
)

func ConnectPostgresDB(db_url string, maxIdleConns, maxOpenConns int, connMaxIdleTime time.Duration) (*sql.DB, error) {

	db, err := sql.Open("postgres", db_url)

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
