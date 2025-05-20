package storage

import (
	"database/sql"

	"golang.org/x/net/context"
)

type AttendeeModel struct {
	db *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	EventID int `json:"event_id"`
}

func (a *AttendeeModel) CreateAttendee(ctx context.Context, attendee *Attendee) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `INSERT INTO attendees (user_id, event_id) VALUES ($1, $2) RETURNING id, user_id, event_id`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = tx.QueryRowContext(ctx, query, attendee.UserID, attendee.EventID).Scan(&attendee.ID, &attendee.UserID, &attendee.EventID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (a *AttendeeModel) GetByEventAndAttendee(ctx context.Context, eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT id, user_id, event_id FROM attendees WHERE event_id = $1 AND user_id = $2`

	attendee := &Attendee{}
	err := a.db.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.ID, &attendee.UserID, &attendee.EventID)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, ErrAttendeeNotFound
		}
		return nil, err
	}

	return attendee, nil
}

func (a *AttendeeModel) GetAttendeesByEvent(ctx context.Context, eventId int) (*[]User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT u.id, u.name, u.email FROM users u
				JOIN attendees a ON u.id = a.user_id
				WHERE a.event_id = $1`

	var users []User

	rows, err := a.db.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (a *AttendeeModel) DeleteAttendee(ctx context.Context, eventId, userId int) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `DELETE FROM attendees WHERE event_id = $1 AND user_id = $2`
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	results, err := tx.ExecContext(ctx, query, eventId, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows == 0 {
		tx.Rollback()
		return ErrAttendeeNotFound
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil

}

func (a *AttendeeModel) GetEventsOfAttendee(ctx context.Context, userId int) (*[]Event, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	query := `SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location FROM events e
	JOIN attendees a ON e.id = a.event_id
	WHERE a.user_id = $1`

	var events []Event

	rows, err := a.db.QueryContext(ctx, query, userId)
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
