package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/puremike/event-mgt-api/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required,min=3"`
	Email string `json:"email" binding:"required,email"`
}

// RegisterUser godoc
//
//	@Summary		Register user
//	@Description	Register a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserRequest	true	"User payload"
//	@Success		200		{object}	storage.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/register [post]
func (app *application) registerUser(c *gin.Context) {

	var payload RegisterUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &storage.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
	}

	if err := app.store.Users.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, response)
}

// GetuserById godoc
//
//	@Summary		Get User
//	@Description	Get User by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	storage.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/auth/{id} [get]
func (app *application) getUserById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := app.store.Users.GetUserByID(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, response)
}
