package models

type User struct{
	ID string `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Password string `json:"password"`
	LoyaltyID string `json:"loyaltyId"`
}