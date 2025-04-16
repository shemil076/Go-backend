package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shemil076/loyalty-backend/internal/services"
)


func EarnPoints(c *gin.Context){
	var loyaltyInputs struct {
		ID int `json:"id" binding:"required"`
		AmountInCents int64 `json:"AmountInCents" binding:"required"`
	}

	if err := c.ShouldBind(&loyaltyInputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}


	resp, err := services.GetUserById(loyaltyInputs.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.EarnPoints(resp.LoyaltyID, loyaltyInputs.AmountInCents)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message":"successfully earned points"})
}