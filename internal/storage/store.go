package storage

import (
	"database/sql"
	"errors"
	"time"
)

type Storage struct {
	Users     UserStore
	Events    EventStore
	Attendees AttendeeStore
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:     UserStore{db},
		Events:    EventStore{db},
		Attendees: AttendeeStore{db},
	}
}

var (
	QueryTimeOutDuration = 5 * time.Second
	ErrEventNotFound     = errors.New("event not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrAttendeeNotFound  = errors.New("attendee not found")
)
