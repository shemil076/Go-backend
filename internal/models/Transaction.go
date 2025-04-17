package models

type Transaction struct{
	Date string `json:"date"`
	Description string `json:"description"`
	Points int64 `json:"points"`
}