package services

import (
	"database/sql"
	"fmt"

	"github.com/shemil076/loyalty-backend/internal/database"
	"github.com/shemil076/loyalty-backend/internal/models"
)


func GetUserById(userId string) (models.User, error) {

	var user models.User
	query := `
		SELECT id, phoneNumber FROM users WHERE id = ?
	`
	err := database.DB.QueryRow(query, userId).Scan(&user.ID, &user.PhoneNumber)

	if (err == sql.ErrNoRows){
		return user,  fmt.Errorf("user not found")
	}else if (err != nil){
		return user, fmt.Errorf("failed to query user: %v", err)
	}

	return user, nil
}