package services

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/shemil076/loyalty-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)


func checkUserExists(phoneNumber string, DB *sql.DB) (bool, error) {
	query := `
		SELECT id FROM users WHERE phoneNumber = ?
	`
	row := DB.QueryRow(query, phoneNumber)

	var id int
	err := row.Scan(&id)

	if (err != nil){
		if (err == sql.ErrNoRows){
			return false, nil
		}
		return false, fmt.Errorf("error checking user: %v", err)
	}
	return true, nil
}

func CreateUser (phoneNumber string, firstName string, lastName string, password string, DB *sql.DB) (models.User, error){
	var newUser models.User

	id := uuid.NewString();
	query := `
		INSERT INTO users (id, phoneNumber, firstName, lastName, password, loyaltyId) VALUES (?, ?, ?, ?, ?, ?)
	`
	userExists, err := checkUserExists(phoneNumber, DB)

	if (err != nil){
		return newUser, err
	}

	if (userExists){
		return newUser, fmt.Errorf("user is already exists for the phone number %v", phoneNumber)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	if (err != nil){
		return newUser, fmt.Errorf("error occurred while hashing password: %v", err)
	}

	loyaltyId, loyaltyError := CreateLoyalAccount(phoneNumber);

	if loyaltyError != nil {
		return newUser, fmt.Errorf("error occurred while creating loyalty account: %v", loyaltyError)
	}

	_, err = DB.Exec(query, id, phoneNumber, firstName, lastName, passwordHash, loyaltyId)

	if (err != nil){
        return newUser, fmt.Errorf("error occurred while creating the user: %v", err)
	}

	// if err != nil {
    //     return newUser, fmt.Errorf("error getting last insert ID: %v", err)
    // }
	newUser.ID = id
	newUser.Password = string(passwordHash)
	newUser.PhoneNumber = phoneNumber
	newUser.LoyaltyID = loyaltyId
	newUser.FirstName = firstName
	newUser.LastName = lastName

	return newUser, nil
} 

func Login(phoneNumber string, password string, DB *sql.DB) (string, error) {
	var user models.User
	var token string 

	query := `
		SELECT id, phoneNumber, password, firstName, lastName, loyaltyId FROM users WHERE phoneNumber = ?
	`


	err := DB.QueryRow(query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.FirstName, &user.LastName, &user.LoyaltyID)

	if (err == sql.ErrNoRows){
		return "", fmt.Errorf("invalid credentials")
	}else if (err != nil){
		return "", fmt.Errorf("failed to query user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) ; err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id" : user.ID,
		"exp" : time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("SECRET")

	if (secret == ""){
		return "", fmt.Errorf("JWT secret key is not configured")
	}

	token, err = generatedToken.SignedString([]byte(secret))
	
	if (err != nil){
		return token, fmt.Errorf("error occurred")
	}

	return token, nil
}