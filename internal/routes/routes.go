package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/handlers"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/user", handlers.CreateUser)
		api.POST("/wallet", handlers.CreateWallet)
		api.GET("/wallet/:user_id/balance", handlers.GetWalletBalance)
		api.POST("/transfer", handlers.TransferMoney)
	}

	return router
}
