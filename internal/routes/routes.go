package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tegveer-singh123/wallet-api/internal/handlers"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/users", handlers.CreateUser)
	}

	return router
}
