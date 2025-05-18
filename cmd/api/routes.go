package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/event-mgt-api/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
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
			events.GET("/:id", app.getEventByID)
			events.PUT("/:id", app.updateEvent)
			events.DELETE("/:id", app.deleteEvent)
		}
	}

	return g
}
