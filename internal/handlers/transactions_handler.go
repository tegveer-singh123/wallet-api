package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/config"
	"github.com/tegveer-singh123/wallet-api/internal/models"
)

func GetUserTransactions(c *gin.Context) {
	userID := c.Param("user_id") //route: /users/:user_id/transactions

	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found for user"})
		return
	}

	var transactions []models.Transaction
	if err := config.DB.
		Where("sender_id = ? OR receiver_id = ?", wallet.UserID, wallet.UserID).
		Order("created_at desc").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	// Add sent/received label
	var enriched []gin.H
	for _, tx := range transactions {
		txType := "received"
		if tx.SenderID == wallet.UserID {
			txType = "sent"
		}
		enriched = append(enriched, gin.H{
			"id":          tx.ID,
			"amount":      tx.Amount,
			"from_wallet": tx.SenderID,
			"to_wallet":   tx.ReceiverID,
			"type":        txType,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      wallet.UserID,
		"wallet_id":    wallet.ID,
		"transactions": enriched,
	})
}
