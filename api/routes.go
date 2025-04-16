package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/api/handlers"
)

func SetupRoutes(router *gin.Engine){
	api := router.Group("/api")
	{
		api.POST("/signup", handlers.CreateUserHandler)
		api.POST("/login", handlers.LoginHandler)
	}
}