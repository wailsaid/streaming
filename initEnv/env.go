package initEnv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load .env file")
		log.Fatal(err)
	}

	if db_host := os.Getenv("db_host");db_host == "" {
		os.Setenv("db_host","127.0.0.1")
	}
	if db_user := os.Getenv("db_user");db_user == "" {
		os.Setenv("db_user","root")
	}
	if db_pwd := os.Getenv("db_password");db_pwd == "" {
		os.Setenv("db_password","root")
	}


}
