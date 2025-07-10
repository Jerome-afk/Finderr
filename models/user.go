package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique;not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	Password     string `json:"password_hash" gorm:"not null"`
	ProfileImage string `json:"profile_image" gorm:"default:'default.png'"`
	SessionToken string `json:"session_token" gorm:"size:255;index"`
}

type AuthForm struct {
	Username string `form:"username"`
	Email	 string `form:"email"`
	Password string `form:"password"`
}