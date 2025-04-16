package models

type User struct{
	ID int `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	Password string `json:"password"`
}