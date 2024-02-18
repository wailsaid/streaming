package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	User     string `form:"user" json:"user"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}
