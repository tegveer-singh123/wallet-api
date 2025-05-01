package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/config"
	"github.com/tegveer-singh123/wallet-api/internal/models"
)

type CreateWalletInput struct {
	UserID uint `form:"user_id" binding:"required"`
}

//this is the create wallet handler function which creates a new wallet for a user
func CreateWallet(c *gin.Context) {
	var input CreateWalletInput

	// Bind the input JSON to the struct
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Check if the user exists
	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

		// Check if the user exists
		var existingWallet models.Wallet
		if err := config.DB.Where("user_id = ?", input.UserID).First(&existingWallet).Error; err == nil {
			// Email already exists, return an error message
			c.JSON(http.StatusConflict, gin.H{
				"error": "Wallet already exists for this user",
			})
			return
		}

	// Create a wallet for the user
	wallet := models.Wallet{
		UserID:  input.UserID,
		Balance: 100, // Initial balance is 100
	}

	// Save the wallet in the database
	if err := config.DB.Create(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create wallet"})
		return
	}

	// Return success response with wallet details
	c.JSON(http.StatusCreated, gin.H{
		"message": "Wallet created successfully",
		"wallet": wallet,
	})
}



//this is the get wallet balance handler function which returns the balance of a wallet
func GetWalletBalance(c *gin.Context) {
	// Get the user_id from the URL parameter
	userID := c.Param("user_id")

	// Find the wallet associated with the user
	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Return the wallet balance
	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet balance fetched successfully",
		"balance": wallet.Balance,
	})
}
