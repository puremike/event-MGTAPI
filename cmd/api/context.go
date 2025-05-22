package main

import (
	"github.com/gin-gonic/gin"
	"github.com/puremike/event-mgt-api/internal/storage"
)

func (app *application) getUserFromContext(c *gin.Context) *storage.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return &storage.User{}
	}
	user, ok := contextUser.(*storage.User)
	if !ok {
		return &storage.User{}
	}
	return user
}

func (app *application) getEventFromContext(c *gin.Context) *storage.Event {
	contextEvent, exists := c.Get("event")
	if !exists {
		return &storage.Event{}
	}

	event, ok := contextEvent.(*storage.Event)
	if !ok {
		return &storage.Event{}
	}
	return event
}
