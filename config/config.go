package config

import (
	"log"

	"github.com/joho/godotenv"
)



func LoadConfig(){
	err := godotenv.Load()

	if (err != nil){
		log.Fatal("No .env file found (loading from system env)")
	}
}