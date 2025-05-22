package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/event-mgt-api/internal/storage"
)

func (app *application) BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is deformed"})
			c.Abort()
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid base64 encoding"})
			c.Abort()
			return
		}

		username := app.config.authConfig.username
		password := app.config.authConfig.password

		decodeStr := strings.TrimSpace(string(decoded))

		creds := strings.SplitN(decodeStr, ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (app *application) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is deformed"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])

		jwtToken, err := app.JWTAuthenticator.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		if jwtToken == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userId, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sub claim type"})
			c.Abort()
			return
		}

		user, err := app.store.Users.GetUserByID(c.Request.Context(), int(userId))
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", user.ID)
		c.Next()
	}
}

func (app *application) eventContextMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		eventId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
			c.Abort()
			return
		}

		event, err := app.store.Events.GetEventByID(c.Request.Context(), eventId)

		if err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
			c.Abort()
			return
		}
		c.Set("event", event)
		c.Next()
	}
}
