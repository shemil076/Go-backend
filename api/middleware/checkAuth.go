package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shemil076/loyalty-backend/internal/services"
)

func CheckAuth(c *gin.Context){
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error" : "Authorization header is missing"})
		c.Abort()
		return
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.Abort()
		return
	}

	tokenString := parts[1]
	
	secret := os.Getenv("SECRET")
	if (secret == ""){
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Server configuration error"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error" : "Invalid or expired token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	userID, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H {"error" : "token expired"})
		c.Abort()
		return
	}

	currentUser,  err := services.GetUserById(userID)

	if (err != nil){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user not found"})
		c.Abort()
		return
	}
	
	c.Set("currentUser", currentUser)
	c.Next()
}