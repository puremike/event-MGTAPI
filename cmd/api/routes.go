package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := g.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "OK", "env", app.config.env, "message", "Health check successful")
		})

		events := v1.Group("/events")
		{
			events.POST("/", app.createEvent)
			events.GET("/", app.getAllEvents)
			events.GET("/:id", app.getEventById)
			events.PUT("/:id", app.updateEvent)
			events.DELETE("/:id", app.deleteEvent)
			events.POST("/:id/attendees/:userId", app.addAttendeeToEvent)
			events.DELETE("/:id/attendees/:userId", app.deleteAttendeeFromEvent)
			events.GET("/:id/attendees", app.getEventAttendees)
		}

		users := v1.Group("/auth")
		{
			users.POST("/register", app.registerUser)
			users.GET("/:id", app.getUserById)
		}

		attendees := v1.Group("/attendees")
		{
			attendees.GET("/:userId/events", app.getEventsOfAttendee)
		}
	}

	return g
}
