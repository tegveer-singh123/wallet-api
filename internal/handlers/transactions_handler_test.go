package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tegveer-singh123/wallet-api/internal/models"
)

func TestGetUserTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Clean up before test
	require.NoError(t, db.Exec("DELETE FROM transactions").Error)
	require.NoError(t, db.Exec("DELETE FROM wallets").Error)
	require.NoError(t, db.Exec("DELETE FROM users").Error)

	// Create users
	userA := models.User{Name: "User A", Email: "usera@example.com"}
	userB := models.User{Name: "User B", Email: "userb@example.com"}
	require.NoError(t, db.Create(&userA).Error)
	require.NoError(t, db.Create(&userB).Error)

	// Create wallets
	walletA := models.Wallet{UserID: userA.ID, Balance: 1000}
	walletB := models.Wallet{UserID: userB.ID, Balance: 500}
	require.NoError(t, db.Create(&walletA).Error)
	require.NoError(t, db.Create(&walletB).Error)

	// Create a transaction from A to B
	tx := models.Transaction{
		SenderID:   userA.ID,
		ReceiverID: userB.ID,
		Amount:     200,
	}
	require.NoError(t, db.Create(&tx).Error)

	router := gin.Default()
	router.GET("/users/:user_id/transactions", GetUserTransactions)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/transactions", userA.ID), nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"type":"sent"`)
	require.Contains(t, recorder.Body.String(), `"from_wallet":`)
	require.Contains(t, recorder.Body.String(), `"to_wallet":`)
	require.Contains(t, recorder.Body.String(), `"amount":200`)
}
