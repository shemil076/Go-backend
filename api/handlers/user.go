package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/internal/services"
)

func GetUser(c *gin.Context) {
	userIdRaw, exists := c.Get("currentUserID")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userId, ok := userIdRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}
	currentUser, err := services.GetUserById(userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          currentUser.ID,
		"phoneNumber": currentUser.PhoneNumber,
		"loyaltyId":   currentUser.LoyaltyID,
		"firstName":   currentUser.FirstName,
		"lastName":    currentUser.LastName,
	})

}
