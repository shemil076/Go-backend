package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/api"
	"github.com/shemil076/loyalty-backend/config"
	"github.com/shemil076/loyalty-backend/internal/database"
)

func main() {
	config.LoadConfig()
	database.InitDB("loyalty.db")
	router := gin.Default()
	api.SetupRoutes(router)
    router.Run(":8080")
}
