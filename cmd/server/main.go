package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()
    router.GET("/ping",getPing)
	router.POST("/postping", postPing)

    router.Run(":8080")
}

func getPing(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func postPing(c *gin.Context){
	c.JSON(200, gin.H{"message": "pong"})
}