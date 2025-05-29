package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const logDir = "log"
const logFilePath = "log/latest"

func InitLogger() {
	// Ensure directory exists
	if _, err := os.Stat(logDir);
	os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	// Rename existing latest log file
	if _, err := os.Stat(logFilePath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		os.Rename(logFilePath, filepath.Join(logDir, "log-"+timestamp))
	}

	// create new latest log
	file, err := os.Create(logFilePath)
	if err != nil {
		fmt.Printf("Could not create log file: %v\n", err)
		return
	}
	file.Close()
}

func LogEvent(user, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s][%s]: %s\n", timestamp, user, message)

	// Append
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not write to log file: %v\n", err)
		return
	}
	defer file.Close()
	file.WriteString(logMessage)
}