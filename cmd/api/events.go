package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/event-mgt-api/internal/storage"
)

type createEventRequest struct {
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        time.Time `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string    `json:"location" binding:"required,min=3"`
}

func (app *application) createEvent(c *gin.Context) {

	var payload createEventRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &storage.Event{
		OwnerID:     c.GetString("userId"),
		Name:        payload.Name,
		Description: payload.Description,
		Date:        payload.Date,
		Location:    payload.Location,
	}

	createdEvent, err := app.store.Events.CreateEvent(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
}
func (app *application) getEventByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID"})
		return
	}
	event, err := app.store.Events.GetEventByID(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, event)
}
func (app *application) getAllEvents(c *gin.Context) {

	events, err := app.store.Events.GetAllEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (app *application) updateEvent(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID"})
		return
	}

	var payload createEventRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &storage.Event{
		OwnerID:     c.GetString("userId"),
		Name:        payload.Name,
		Description: payload.Description,
		Date:        payload.Date,
		Location:    payload.Location,
	}

	updatedEvent, err := app.store.Events.UpdateEvent(c.Request.Context(), event, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}
func (app *application) deleteEvent(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID"})
		return
	}

	if err := app.store.Events.DeleteEvent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "event deleted successfully"})
}
