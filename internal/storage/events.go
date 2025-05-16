package storage

import (
	"database/sql"
	"time"
)

type EventModel struct {
	db *sql.DB
}

type Events struct {
	ID          int       `json:"id"`
	OwnerID     string    `json:"ownerId" binding:"required"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        time.Time `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string    `json:"location" binding:"required,min=3"`
}
