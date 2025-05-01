package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/config"
	"github.com/tegveer-singh123/wallet-api/internal/models"
)

type CreateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

//this is the create user handler function which creates a new user
func CreateUser(c *gin.Context) {
	var input CreateUserInput

	// Bind the input JSON to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Check if the user exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		// Email already exists, return an error message
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email is already taken. Please use a different email address.",
		})
		return
	}

	user := models.User{Name: input.Name, Email: input.Email}

	//save the user in DB
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": user,
	})
}
