package main

import (
	"context"
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

		jwtToken, err := app.jWTAuthenticator.ValidateToken(token)
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

		user, err := app.getUserFromCache(c.Request.Context(), int(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", user.ID)
		c.Next()
	}
}

func (app *application) getUserFromCache(ctx context.Context, id int) (*storage.User, error) {

	if !app.config.redisClientConfig.enabled {
		return app.store.Users.GetUserByID(ctx, id)
	}
	app.logger.Infow("cache hit", "key", "id", id)

	user, err := app.cacheStorage.Users.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("fetching user from DB", "key", "id", id)
		user, err := app.store.Users.GetUserByID(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				return nil, storage.ErrUserNotFound
			}
			return nil, err
		}
		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			app.logger.Errorw("failed to set user in cache", "key", "id", id, "error", err)
		}
		return user, nil
	}
	return user, nil
}

func (app *application) eventContextMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		eventId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
			c.Abort()
			return
		}

		event, err := app.getEventFromCache(c.Request.Context(), eventId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
			c.Abort()
			return
		}
		c.Set("event", event)
		c.Next()
	}
}

func (app *application) getEventFromCache(ctx context.Context, id int) (*storage.Event, error) {

	if !app.config.redisClientConfig.enabled {
		return app.store.Events.GetEventByID(ctx, id)
	}
	app.logger.Infow("cache hit", "key", "id", id)

	event, err := app.cacheStorage.Events.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if event == nil {
		app.logger.Infow("fetching event from DB", "key", "id", id)
		event, err := app.store.Events.GetEventByID(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				return nil, storage.ErrEventNotFound
			}
			return nil, err
		}
		if err := app.cacheStorage.Events.Set(ctx, event); err != nil {
			app.logger.Errorw("failed to set event in cache", "key", "id", id, "error", err)
		}
		return event, nil
	}
	return event, nil
}
