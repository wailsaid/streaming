package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/saidwail/streaming/models"
)

var (
	DB     *gorm.DB
	config struct {
		db_conn string
		db_name string
		db_host string
		db_user string
		db_pwd  string
		db_port string
	}
)

func Init() {
	if config.db_conn = os.Getenv("DB_CONNECT"); config.db_conn == "" {
		config.db_conn = "sqlite"
	}
	if config.db_host = os.Getenv("DB_HOST"); config.db_host == "" {
		config.db_host = "127.0.1.1"
	}
	if config.db_name = os.Getenv("DB_NAME"); config.db_name == "" {
		config.db_name = "db"
	}

	if config.db_port = os.Getenv("DB_PORT"); config.db_port == "" {
		config.db_port = "3306"
	}

	if config.db_user = os.Getenv("DB_USER"); config.db_user == "" {
		config.db_user = "root"
	}

	if config.db_pwd = os.Getenv("DB_PASSWORD"); config.db_pwd == "" {
		config.db_pwd = ""
	}

}

func Connect() {

	db, err := getConnection(config.db_conn)
	DB = db

	if err != nil {
		log.Fatalf("could not connect to a database: %s", err.Error())
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("DB: could not create a user table")
	}
	log.Printf("user table done")

	err = DB.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatal("DB: could not create a video table")
	}
	log.Printf("video table done")

	err = DB.AutoMigrate(&models.Comment{})
	if err != nil {
		log.Fatal("DB: could not create a comment table")
	}
	log.Printf("comment table done")

	fmt.Println("DB: Done ok")
}

func getConnection(connection string) (*gorm.DB, error) {
	switch connection {
	case "mysql":
		return mysqlConnect()

	case "postgres":
		return postgresConnect()

	case "sqlite":
		return gorm.Open(sqlite.Open("database.sql"), &gorm.Config{})
	}

	return nil, errors.New("unsupport db connection")
}

func mysqlConnect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.db_user, config.db_pwd, config.db_host, config.db_port, config.db_name)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})

}

func postgresConnect() (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.db_host, config.db_user, config.db_pwd, config.db_name, config.db_port)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
