package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/saidwail/streaming/models"
)

var DB *gorm.DB
var config struct {
	db_name string
	db_host string
	db_user string
	db_pwd  string
	db_port string
}

func Connect() {

	//var err error
	dbpath := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.db_user, config.db_pwd, config.db_host, config.db_port, config.db_name)

	log.Printf("URL : %s", dbpath)

	DB, err := gorm.Open(mysql.Open(dbpath), &gorm.Config{})

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
	log.Printf("user table done")

	fmt.Println("DB: Done ok")
}

func Init() {

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

func GetAllVideos() []models.Video {
	if DB == nil {
		log.Println("DB is not initialized")
		return nil
	}
	log.Println("DB is initialized")

	var videos []models.Video
	DB.Find(&videos)
	return videos
}

// New methods to add:

func CreateVideo(video *models.Video) error {
	return DB.Create(video).Error
}

func FindVideoByID(id string) (models.Video, error) {
	var video models.Video
	err := DB.First(&video, id).Error
	return video, err
}

func FindAllVideos() ([]models.Video, error) {
	var videos []models.Video
	err := DB.Find(&videos).Error
	return videos, err
}

func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

func FindUserByEmail(email string) (models.User, error) {
	var user models.User
	err := DB.First(&user, "email = ?", email).Error
	return user, err
}

func FindAllUsers() ([]models.User, error) {
	var users []models.User
	err := DB.Find(&users).Error
	return users, err
}
