package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tegveer-singh123/wallet-api/internal/config"
	"github.com/tegveer-singh123/wallet-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	err := godotenv.Load("./../../.env.test")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dsn := os.Getenv("TEST_DB_DSN") // e.g. "host=localhost user=postgres password=postgres dbname=wallet_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Migrate the models
	err = db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
	require.NoError(t, err)

	config.DB = db
	return db
}

// TestCreateUser is a unit test for the CreateUser handler.
// It tests that CreateUser correctly inserts a new user into the
// users table and returns a success message.
func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Clean up users table before test
	require.NoError(t, db.Exec("DELETE FROM users").Error)

	router := gin.Default()
	router.POST("/users", CreateUser)

	payload := CreateUserInput{
		Name:  "Tegveer Test",
		Email: "tegveer@example.com",
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Contains(t, recorder.Body.String(), "User created successfully")
}
