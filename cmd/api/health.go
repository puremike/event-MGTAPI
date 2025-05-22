package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetHealth godoc
//
//	@Summary		Health Check
//	@Description	Returns the health status of the application
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Health check successful"
//	@Failure		500	{object}	error
//	@Router			/health [get]
//	@Security		BasicAuth
func (app *application) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"env":     app.config.env,
		"message": "Health check successful",
		"time":    time.Now().Format(time.RFC3339),
	})
}
