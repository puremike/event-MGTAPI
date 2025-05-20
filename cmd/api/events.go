package main

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/event-mgt-api/internal/storage"
)

type CreateEventRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3"`
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// CreateEvent godoc
//
//	@Summary		Create event
//	@Description	Create a new event
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateEventRequest	true	"Event payload"
//	@Success		200		{object}	storage.Event
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/events [post]
func (app *application) createEvent(c *gin.Context) {

	var payload CreateEventRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. expected YYYY-MM-DD"})
		return
	}

	event := &storage.Event{
		OwnerID:     r.Intn(3),
		Name:        payload.Name,
		Description: payload.Description,
		Date:        date,
		Location:    payload.Location,
	}

	if err := app.store.Events.CreateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEventByID godoc
//
//	@Summary		Get Event
//	@Description	Get Event by ID
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	storage.Event
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/events/{id} [get]
func (app *application) getEventById(c *gin.Context) {

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID"})
		return
	}
	event, err := app.store.Events.GetEventByID(c.Request.Context(), eventId)

	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		app.logger.Errorf("failed to retrieve event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetEvents godoc
//
//	@Summary		Get Events
//	@Description	Get All Events
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	storage.Event
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/events [get]
func (app *application) getAllEvents(c *gin.Context) {

	events, err := app.store.Events.GetAllEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// UpdateEvent godoc
//
//	@Summary		Update event
//	@Description	Update event by ID
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateEventRequest	true	"Event payload"
//	@Param			id		path		int					true	"Event ID"
//
//	@Success		200		{object}	storage.Event
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/events/{id} [put]
func (app *application) updateEvent(c *gin.Context) {

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID"})
		return
	}

	var payload CreateEventRequest

	date, _ := time.Parse("2006-01-02", payload.Date)

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &storage.Event{
		OwnerID:     r.Intn(500),
		Name:        payload.Name,
		Description: payload.Description,
		Date:        date,
		Location:    payload.Location,
	}

	updatedEvent, err := app.store.Events.UpdateEvent(c.Request.Context(), event, eventId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent godoc
//
//	@Summary		Delete event
//	@Description	Delete event by ID
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Event ID"
//
//	@Success		204	{string}	string	"no content"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context) {

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID"})
		return
	}

	if err := app.store.Events.DeleteEvent(c.Request.Context(), eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event"})
		return
	}

	c.Status(http.StatusNoContent)
}

// addAttendeeToEvent adds a user as an attendee to a specific event.
//
//	@Summary		Add an attendee to an event
//	@Description	Adds a user to the list of attendees for a given event by event ID and user ID.
//	@Tags			Attendees
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Event ID"
//	@Param			userId	path		int					true	"User ID"
//	@Success		201		{object}	storage.Attendee	"Attendee successfully added"
//	@Failure		400		{object}	map[string]string	"Invalid event ID or user ID"
//	@Failure		404		{object}	map[string]string	"Event or user not found"
//	@Failure		409		{object}	map[string]string	"Attendee already exists"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/events/{id}/attendees/{userId} [post]
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	event, err := app.store.Events.GetEventByID(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}

	user, err := app.store.Users.GetUserByID(c.Request.Context(), userId)
	if err != nil {
		app.logger.Errorf("failed to retrieve attendee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	existingAttendee, err := app.store.Attendees.GetByEventAndAttendee(c.Request.Context(), event.ID, user.ID)
	if err != nil && !errors.Is(err, storage.ErrAttendeeNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "attendee already exists"})
		return
	}

	attendee := &storage.Attendee{
		UserID:  user.ID,
		EventID: event.ID,
	}

	if err := app.store.Attendees.CreateAttendee(c.Request.Context(), attendee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attendee"})
		return
	}
	c.JSON(http.StatusCreated, attendee)
}

// GetEventAttendees get the attendees to a specific event.
//
//	@Summary		Get event attendees
//	@Description	Get the list of attendees for a given event by event ID.
//	@Tags			Attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Event ID"
//	@Success		200	{object}	storage.Attendee	"Attendees successfully retrieved"
//	@Failure		400	{object}	map[string]string	"Invalid event ID"
//	@Failure		404	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/events/{id}/attendees [get]
func (app *application) getEventAttendees(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := app.store.Attendees.GetAttendeesByEvent(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event attendees"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteAttendee godoc
//
//	@Summary		Delete attendee
//	@Description	Delete attendee by event and user ID
//	@Tags			Attendees
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Event ID"
//	@Param			userId	path		int		true	"User ID"
//	@Success		204		{string}	string	"no content"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/events/{id}/attendees/{userId} [delete]
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	attendee, err := app.store.Attendees.GetByEventAndAttendee(c.Request.Context(), eventId, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAttendeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "attendee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendee"})
		return
	}

	if err := app.store.Attendees.DeleteAttendee(c.Request.Context(), attendee.EventID, attendee.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attendee"})
		return
	}

	c.Status(http.StatusNoContent)
}
