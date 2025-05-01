package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tegveer-singh123/wallet-api/internal/models"
)

func TestTransferMoney(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Clean tables before test
	require.NoError(t, db.Exec("DELETE FROM transactions").Error)
	require.NoError(t, db.Exec("DELETE FROM wallets").Error)
	require.NoError(t, db.Exec("DELETE FROM users").Error)

	// Create sender and receiver users
	sender := models.User{Name: "Sender User", Email: "sender@example.com"}
	receiver := models.User{Name: "Receiver User", Email: "receiver@example.com"}
	require.NoError(t, db.Create(&sender).Error)
	require.NoError(t, db.Create(&receiver).Error)

	// Create wallets with balances
	senderWallet := models.Wallet{UserID: sender.ID, Balance: 500.0}
	receiverWallet := models.Wallet{UserID: receiver.ID, Balance: 100.0}
	require.NoError(t, db.Create(&senderWallet).Error)
	require.NoError(t, db.Create(&receiverWallet).Error)

	router := gin.Default()
	router.POST("/transfer", TransferMoney)

	// Transfer payload
	payload := TransferInput{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Amount:     150.0,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), "Transfer successful")

	// Check if balances updated correctly
	var updatedSender, updatedReceiver models.Wallet
	require.NoError(t, db.Where("user_id = ?", sender.ID).First(&updatedSender).Error)
	require.NoError(t, db.Where("user_id = ?", receiver.ID).First(&updatedReceiver).Error)

	require.Equal(t, 350.0, updatedSender.Balance)
	require.Equal(t, 250.0, updatedReceiver.Balance)

	// Check if transaction was recorded
	var txn models.Transaction
	require.NoError(t, db.First(&txn).Error)
	require.Equal(t, sender.ID, txn.SenderID)
	require.Equal(t, receiver.ID, txn.ReceiverID)
	require.Equal(t, payload.Amount, txn.Amount)
}
