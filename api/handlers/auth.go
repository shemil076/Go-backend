package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/internal/database"
	"github.com/shemil076/loyalty-backend/internal/models"
	"github.com/shemil076/loyalty-backend/internal/services"
)

func CreateUserHandler(c *gin.Context){
	var cred models.AuthInput

	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := services.CreateUser(cred.PhoneNumber, cred.Password, database.DB)

	if (err != nil){
		if (strings.Contains(err.Error(), "user already exists")){
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
		"phoneNumber": user.PhoneNumber,
	});
}

func LoginHandler(c *gin.Context){
	var cred models.AuthInput

	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	token, err := services.Login(cred.PhoneNumber, cred.Password, database.DB)

	if (err != nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}