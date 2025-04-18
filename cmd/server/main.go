package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/api"
	"github.com/shemil076/loyalty-backend/config"
	"github.com/shemil076/loyalty-backend/internal/database"
)

func main() {
	config.LoadConfig()
	database.InitDB("loyalty.db")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api.SetupRoutes(router)
    router.Run(":8080")
}
