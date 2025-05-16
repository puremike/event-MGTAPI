package storage

import "database/sql"

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
