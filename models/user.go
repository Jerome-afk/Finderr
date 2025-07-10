package models

import (
	"gorm.io/gorm"
)

type Role string

const (
	AdminRole    Role = "admin"
	MemberRole   Role = "member"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique;not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	Password     string `json:"password_hash" gorm:"not null"`
	Role         Role   `json:"role" gorm:"default:'member'"`
	ProfileImage string `json:"profile_image" gorm:"default:'/images/default.jpg'"`
	SessionToken string `json:"session_token" gorm:"size:255;index"`
}

type AuthForm struct {
	Username string `form:"username"`
	Email	 string `form:"email"`
	Password string `form:"password"`
}