package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/api/handlers"
	"github.com/shemil076/loyalty-backend/api/middleware"
)

func SetupRoutes(router *gin.Engine){
	api := router.Group("/api")
	{
		api.POST("/signup", handlers.CreateUserHandler)
		api.POST("/login", handlers.LoginHandler)
		api.POST("/earn", handlers.EarnPoints)
		api.POST("/redeem", handlers.RedeemLoyaltyReward)
		api.GET("/balance", handlers.GetBalance)
		api.GET("/history", handlers.GetTransactions)
		api.GET("/user", middleware.CheckAuth(),handlers.GetUser)
	}
}