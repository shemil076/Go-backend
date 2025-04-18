package models

type Transaction struct{
	Date string `json:"date"`
	Title string `json:"title"`
	Description string `json:"description"`
	Points int64 `json:"points"`
}