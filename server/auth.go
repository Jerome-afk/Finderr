package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

// Checks if username and password match
func Authenticate(username, password string) bool {
	filePath := fmt.Sprintf("users/%s.json", username)

	// Check is user file exists
	file, err := os.Open(filePath)
	if err != nil {
		logEvent("server", "Unauthorized login attempt")
		return false
	}
	defer file.Close()

	// Parse JSON file
	var userData struct {
		HashPassword string `json:"hash_password"`
	}
	if err := json.NewDecoder(file).Decode(&userData); err != nil {
		logEvent("server", "Failed to get user password")
		return false
	}

	// Hash the entered password
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	if hashedPassword != userData.HashPassword {
		logEvent(username, "Wrong password by user")
		return false
	}
	logEvent(username, "User logged in")
	return true
}