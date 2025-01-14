package database

import (
	"log"

	"github.com/saidwail/streaming/models"
)

func GetAllVideos() []models.Video {
	if DB == nil {
		log.Println("DB is not initialized")
		return nil
	}

	var videos []models.Video
	DB.Find(&videos)
	return videos
}

// New methods to add:

func CreateVideo(video *models.Video) error {
	res := DB.Create(video)
	if res.Error != nil {
		log.Printf("Error creating video: %v", res.Error.Error())
		return res.Error
	}
	return nil
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

func GetRecommendedVideos(currentVideoID string, limit int) []models.Video {
	var videos []models.Video
	DB.Where("id != ?", currentVideoID).Limit(limit).Find(&videos)
	return videos
}
