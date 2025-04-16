package models

type AuthInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password string `json:"password" binding:"required"`
}