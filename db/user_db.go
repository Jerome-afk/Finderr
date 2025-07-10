package db

import (
	"crypto/rand"
	"encoding/base64"

	"finderr/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func SetUserSession(userID uint, token string) error {
	var user models.User
	if err := DB.Model(&user).Where("id = ?", userID).Update("session_token", token).Error; err != nil {
		return err
	}
	return nil
}

func GetUserBySessionToken(token string) (*models.User, error) {
	var user models.User
	if err := DB.Where("session_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserSession(userID uint) error {
	var user models.User
	if err := DB.Model(&user).Where("id = ?", userID).Update("session_token", "").Error; err != nil {
		return err
	}
	return nil
}
