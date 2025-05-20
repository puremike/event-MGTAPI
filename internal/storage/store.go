package storage

import (
	"database/sql"
	"errors"
	"time"
)

type Storage struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeeModel
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:     UserModel{db},
		Events:    EventModel{db},
		Attendees: AttendeeModel{db},
	}
}

var (
	QueryTimeOutDuration = 5 * time.Second
	ErrEventNotFound     = errors.New("event not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrAttendeeNotFound  = errors.New("attendee not found")
)
