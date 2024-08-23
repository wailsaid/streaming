package database

import (
	"fmt"
	"log"

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
