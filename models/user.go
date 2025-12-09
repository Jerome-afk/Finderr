package models

import (
	"time"
)

type Role string

const (
	AdminRole  Role = "admin"
	MemberRole Role = "member"
)

type User struct {
	ID        uint   `gorm:"primaryKey;not null;unique"`
	Username  string `json:"username" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password_hash" gorm:"not null"`
	Role      Role   `json:"role" gorm:"default:'member'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthForm struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}
