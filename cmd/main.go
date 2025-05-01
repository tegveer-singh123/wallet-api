package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/config"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()

	port := os.Getenv("APP_PORT")

	r := gin.Default()

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Wallet API is running"})
	})

	err := r.Run(":"+port)
	if err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
