package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEventsOfAnAttendee  get the events of an attendee.
//
//	@Summary		Get Attendee events
//	@Description	Get the list of events for a given attendee.
//	@Tags			Attendees
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int					true	"User ID"
//	@Success		200		{object}	storage.Event		"Events successfully retrieved"
//	@Failure		400		{object}	map[string]string	"Invalid user ID"
//	@Failure		409		{object}	error
//	@Failure		500		{object}	error
//	@Router			/attendees/{userId}/events [get]
func (app *application) getEventsOfAttendee(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	events, err := app.store.Attendees.GetEventsOfAttendee(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
