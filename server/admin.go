package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
)

const adminFilePath = "users/admin.json"

// Check and create admin file
func checkAndCreateAdminFile() {
	// Ensure user directory exists
	if _, err := os.Stat("users");
	os.IsNotExist(err) {
		os.Mkdir("users", 0755)
	}

	// Check if it is empty
	files, err := ioutil.ReadDir("users")
	if err != nil || len(files) > 0 {
		return
	}

	// Generate random password
	password := generateRandomPassword(8)

	// Hash the password
	hashedPassword := hashPassword(password)

	// Create admin file
	adminData := `{
	"hash_password": "`+ hashedPassword + `",
	"user_level": 1,
	"email": "nobody@email.com"
	}`

	// Write admin data to file
	err = ioutil.WriteFile(adminFilePath, []byte(adminData), 0644)
	if err != nil {
		logEvent("server", "Failed to create admin file:" + err.Error())
		return
	}

	os.Mkdir("admin", 0755)
	logEvent("server", "user admin created successfully")
}

// generateRandomPassword creates a random password of the given length
func generateRandomPassword(length int) string {
	// Generate random bytes for password
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		logEvent("server", "Failed to generate random password: "+err.Error())
		return ""
	}

	// Convert random bytes to hex string and truncate to the required length
	password := hex.EncodeToString(bytes)[:length]

	// Log the generated raw password (before encoding)
	logEvent("server", "Generated raw password: "+password)

	return password
}

// hashPassword hashes the password using SHA-256
func hashPassword(password string) string {
	// Hash the password using SHA-256
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)

	// Encode the hashed password into a hex string
	return hex.EncodeToString(hashedPassword)
}