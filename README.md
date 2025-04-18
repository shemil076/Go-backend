Square Loyalty API Backend
Setup Instructions

Clone the Repository:
git clone https://github.com/shemil076/Go-backend.git
cd loyalty-backend


Install Dependencies:
go mod tidy


Set Up SQLite (Optional): If using SQLite, initialize the database:
sqlite3 users.db < schema.sql



How to Run the Application

Start the Server: Run the Go application:
go run cmd/server/main.go 

The server starts on http://localhost:8080.



