package storage

import (
	"context"
	"database/sql"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	OwnerID     int       `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}

type EventStore struct {
	db *sql.DB
}

func (e *EventStore) CreateEvent(ctx context.Context, event *Event) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING owner_id, name, description, date, location`

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, query, event.OwnerID, event.Name, event.Description, event.Date, event.Location).Scan(&event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (e *EventStore) GetEventByID(ctx context.Context, eventId int) (*Event, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events WHERE id = $1`

	event := &Event{}

	err := e.db.QueryRowContext(ctx, query, eventId).Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (e *EventStore) GetAllEvents(ctx context.Context) (*[]Event, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events`

	var events []Event

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		if err = rows.Scan(&e.ID, &e.OwnerID, &e.Name, &e.Description, &e.Date, &e.Location); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &events, nil
}

func (e *EventStore) UpdateEvent(ctx context.Context, event *Event, eventId int) (*Event, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `UPDATE events SET name = $1, description = $2, date = $3, location = $4 WHERE id = $5 RETURNING owner_id, name, description, date, location`

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, query, event.Name, event.Description, event.Date, event.Location, eventId).Scan(&event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return event, nil
}

func (e *EventStore) DeleteEvent(ctx context.Context, eventId int) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, query, eventId)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows == 0 {
		tx.Rollback()
		return ErrEventNotFound
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
