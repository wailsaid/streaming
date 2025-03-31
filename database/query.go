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

// GetPaginatedVideos returns videos with pagination
func GetPaginatedVideos(offset, limit int) []models.Video {
	var videos []models.Video
	DB.Offset(offset).Limit(limit).Find(&videos)
	return videos
}

// SearchVideos searches videos by title or description
func SearchVideos(query string) []models.Video {
	var videos []models.Video
	searchQuery := "%" + query + "%"
	DB.Where("title LIKE ? OR description LIKE ?", searchQuery, searchQuery).Find(&videos)
	return videos
}

// CreateComment adds a new comment
func CreateComment(comment *models.Comment) error {
	return DB.Create(comment).Error
}

// GetCommentsByVideoID retrieves comments for a video
func GetCommentsByVideoID(videoID uint, limit int) []models.Comment {
	var comments []models.Comment
	DB.Where("video_id = ?", videoID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&comments)
	return comments
}

// FindUserByID retrieves a user by ID
func FindUserByID(id uint) (models.User, error) {
	var user models.User
	err := DB.First(&user, id).Error
	return user, err
}
