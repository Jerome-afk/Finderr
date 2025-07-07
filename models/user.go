package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique;not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	PasswordHash string `json:"password_hash" gorm:"not null"`
	ProfileImage string `json:"profile_image" gorm:"default:'default.png'"`
}