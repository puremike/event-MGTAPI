package storage

import (
	"context"
	"database/sql"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	OwnerID     string    `json:"ownerId" binding:"required"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        time.Time `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string    `json:"location" binding:"required,min=3"`
}

type EventModel struct {
	db *sql.DB
}

func (e *EventModel) CreateEvent(ctx context.Context, event *Event) (*Event, error) {
	return nil, nil
}

func (e *EventModel) GetEventByID(ctx context.Context, id int) (*Event, error) {
	return &Event{}, nil
}

func (e *EventModel) GetAllEvents(ctx context.Context) (*[]Event, error) {
	return nil, nil
}

func (e *EventModel) UpdateEvent(ctx context.Context, event *Event, id int) (*Event, error) {
	return nil, nil
}

func (e *EventModel) DeleteEvent(ctx context.Context, id int) error {
	return nil
}
