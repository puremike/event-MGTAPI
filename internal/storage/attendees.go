package storage

import "database/sql"

type AttendeeModel struct {
	db *sql.DB
}

type Attendees struct {
	ID      int    `json:"id"`
	UserID  string `json:"userId"`
	EventID string `json:"eventId"`
}
