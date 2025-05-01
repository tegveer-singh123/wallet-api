package main

import (
	"log"
	"os"

	"github.com/tegveer-singh123/wallet-api/internal/config"
	"github.com/tegveer-singh123/wallet-api/internal/routes"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()

	port := os.Getenv("APP_PORT")

	r := routes.SetupRoutes()

	err := r.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
