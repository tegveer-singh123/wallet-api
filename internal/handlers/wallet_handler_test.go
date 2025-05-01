package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tegveer-singh123/wallet-api/internal/models"
)

func TestCreateWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Clean up before test
	require.NoError(t, db.Exec("DELETE FROM wallets").Error)
	require.NoError(t, db.Exec("DELETE FROM users").Error)

	// Create a test user
	user := models.User{Name: "Test Wallet User", Email: "walletuser@example.com"}
	require.NoError(t, db.Create(&user).Error)

	router := gin.Default()
	router.POST("/wallet", CreateWallet)

	payload := CreateWalletInput{
		UserID: user.ID,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Contains(t, recorder.Body.String(), "Wallet created successfully")
}

func TestGetWalletBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Clean up before test
	require.NoError(t, db.Exec("DELETE FROM wallets").Error)
	require.NoError(t, db.Exec("DELETE FROM users").Error)

	// Create a test user and wallet
	user := models.User{Name: "Balance Test User", Email: "balance@example.com"}
	require.NoError(t, db.Create(&user).Error)

	wallet := models.Wallet{UserID: user.ID, Balance: 250.75}
	require.NoError(t, db.Create(&wallet).Error)

	router := gin.Default()
	router.GET("/wallets/:user_id/balance", GetWalletBalance)

	req, err := http.NewRequest(http.MethodGet, "/wallets/"+fmt.Sprint(user.ID)+"/balance", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), "Wallet balance fetched successfully")
	require.Contains(t, recorder.Body.String(), "250.75")
}
