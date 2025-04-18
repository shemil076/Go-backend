package models

type AuthInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
}