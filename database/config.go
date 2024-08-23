package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/saidwail/streaming/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.sql"))

	if err != nil {
		log.Fatal("DB: could not connect to a database")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("DB: could not create a user table")
	}
	err = DB.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatal("DB: could not create a video table")
	}
	fmt.Println("DB: Done ok")
}

func Init() {

	if db_host := os.Getenv("db_host"); db_host == "" {
		os.Setenv("db_host", "127.0.0.1")
	}
	if db_user := os.Getenv("db_user"); db_user == "" {
		os.Setenv("db_user", "root")
	}
	if db_pwd := os.Getenv("db_password"); db_pwd == "" {
		os.Setenv("db_password", "root")
	}
}
