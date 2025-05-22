package main

import (
	"expvar"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/puremike/event-mgt-api/internal/env"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	// Add CORS middleware
	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env.GetEnvString("CORS_ALLOWED_ORIGIN", "https://yourfrontend.com")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := g.Group("/api/v1")
	{
		basicAuth := v1.Group("/")
		basicAuth.Use(app.BasicAuthMiddleware())
		{
			basicAuth.GET("/debug/vars", gin.WrapH(app.expvars(expvar.Handler())))
			basicAuth.GET("/health", app.healthCheck)
		}

		events := v1.Group("/events")
		{
			events.GET("/", app.getAllEvents)
			events.GET("/:id", app.getEventById)
			events.GET("/:id/attendees", app.getEventAttendees)
		}

		users := v1.Group("/auth")
		{
			users.POST("/register", app.registerUser)
			users.GET("/:id", app.getUserById)
			users.POST("/login", app.loginUser)
		}

		attendees := v1.Group("/attendees")
		{
			attendees.GET("/:userId/events", app.getEventsOfAttendee)
		}

		authGroup := v1.Group("/")
		authGroup.Use(app.AuthMiddleware())
		{
			authGroup.POST("/events", app.createEvent)
			authGroup.Use(app.eventContextMiddleWare()).PUT("/events/:id", app.updateEvent)
			authGroup.Use(app.eventContextMiddleWare()).DELETE("/events/:id", app.deleteEvent)
			authGroup.Use(app.eventContextMiddleWare()).POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
			authGroup.Use(app.eventContextMiddleWare()).DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEvent)
		}
	}

	return g
}
