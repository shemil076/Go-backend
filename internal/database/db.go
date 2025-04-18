package database

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dsnURI string) {
	var err error

	DB, err = sql.Open("sqlite", dsnURI)

	if err != nil {
		log.Fatalf("Failed to connect to the database %v\n", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v\n", err)
	}

	createTable()
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY ,
		firstName TEXT,
		lastName TEXT,
		phoneNumber TEXT UNIQUE,
		loyaltyId TEXT UNIQUE,
		password TEXT
	);`

	_, err := DB.Exec(query)

	if err != nil {
		log.Fatalf("Failed to create users table %v\n", err)
	}
}
