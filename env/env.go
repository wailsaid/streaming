package env

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/saidwail/streaming/database"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	database.Init()
}
