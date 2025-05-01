package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/config"
	"github.com/tegveer-singh123/wallet-api/internal/models"
	"gorm.io/gorm"
)

type TransferInput struct {
	SenderID   uint    `json:"sender_id" binding:"required"`
	ReceiverID uint    `json:"receiver_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

func TransferMoney(c *gin.Context) {
	var input TransferInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.SenderID == input.ReceiverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to the same wallet"})
		return
	}

	var transaction models.Transaction
	var senderWallet, receiverWallet models.Wallet

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", input.SenderID).First(&senderWallet).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sender wallet not found"})
			return err
		}

		if senderWallet.Balance < input.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return fmt.Errorf("Insufficient balance")
		}

		if err := tx.Where("user_id = ?", input.ReceiverID).First(&receiverWallet).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Receiver wallet not found"})
			return err
		}

		senderWallet.Balance -= input.Amount
		receiverWallet.Balance += input.Amount

		if err := tx.Save(&senderWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&receiverWallet).Error; err != nil {
			return err
		}

		transaction = models.Transaction{
			SenderID:   input.SenderID,
			ReceiverID: input.ReceiverID,
			Amount:     input.Amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	// 🪵 Log full transfer details
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transfer successful",
		"transaction_id": transaction.ID,
		"amount":         transaction.Amount,
		"from_wallet":    transaction.SenderID,
		"to_wallet":      transaction.ReceiverID,
		"balances": gin.H{
			"sender_balance":   senderWallet.Balance,
			"receiver_balance": receiverWallet.Balance,
		},
	})
}
