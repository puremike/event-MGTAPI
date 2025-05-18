package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "OK", "env", app.config.env, "message", "Health check successful")
		})

		events := v1.Group("/events")
		{
			events.POST("/", app.createEvent)
			events.GET("/", app.getAllEvents)
			events.GET("/:id", app.getEventByID)
			events.PUT("/:id", app.updateEvent)
			events.DELETE("/:id", app.deleteEvent)
		}
	}

	return g
}
