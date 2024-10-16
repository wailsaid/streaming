package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	User      string `form:"user" json:"user"`
	Email     string `form:"email" json:"email"`
	Password  string `form:"password" json:"-"`
	ProfileID uint
	Profile   Profile
}

type Profile struct {
	gorm.Model
	PicturePath string
}

type Channel struct {
	Videos []Video
}

type Video struct {
	gorm.Model
	Title         string `form:"title" json:"title"`
	VideoPath     string `form:"video" json:"video"`
	ThumbnailPath string `form:"thumbnail" json:"thumbnail"`
	Description   string `form:"description" json:"description"`
}
