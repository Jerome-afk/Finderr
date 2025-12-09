package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Time      time.Time   `json:"time"`
	Data      interface{} `json:"data,omitempty"`
}

const (
	logDir     = "logs"
	latestFile = "latest"
	maxLogSize = 1024 * 1024 // 1MB threshold
	archiveDir = "archive"
)

var latestFilePath = filepath.Join(logDir, latestFile)

// Ensure folders exist
func initLogFolders() error {
	// First create the log folder
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create logs directory: %w", err)
		}
	}

	// Create an archive folder
	arc := filepath.Join(logDir, archiveDir)
	if _, err := os.Stat(arc); os.IsNotExist(err) {
		if err := os.Mkdir(arc, 0755); err != nil {
			return fmt.Errorf("failed to create archive directory: %w", err)
		}
	}

	return nil
}

func rotateLogIfNeeded() error {
	// Check the file size
	info, err := os.Stat(latestFilePath)
	if os.IsNotExist(err) {
		return nil // No file yet to rotate
	}

	if err != nil {
		return fmt.Errorf("stat error: %w", err)
	}

	if info.Size() < maxLogSize {
		return nil // Still small
	}

	// Rotate if file is bigger than 1MB
	timestamp := time.Now().Format("20060102_150405")
	rotatedName := filepath.Join(logDir, archiveDir, fmt.Sprintf("%s.log.gz", timestamp))

	oldFile, err := os.Open(latestFilePath)
	if err != nil {
		return fmt.Errorf("open latest failed: %w", err)
	}
	defer oldFile.Close()

	// Create a gzip archive
	outFile, err := os.Create(rotatedName)
	if err != nil {
		return fmt.Errorf("create archive failed: %w", err)
	}
	gz := gzip.NewWriter(outFile)

	// Copy data into gzip
	if _, err := io.Copy(gz, oldFile); err != nil {
		gz.Close()
		outFile.Close()
		return fmt.Errorf("gzip copy failed: %w", err)
	}

	gz.Close()
	outFile.Close()

	os.Remove(latestFilePath)

	// Create new
	_, err = os.Create(latestFilePath)
	if err != nil {
		return fmt.Errorf("creating new latest failed: %w", err)
	}

	return nil
}

func LogHandler() error {
	if err := initLogFolders(); err != nil {
		return err
	}

	// Create latest if missing
	if _, err := os.Stat(latestFilePath); os.IsNotExist(err) {
		_, err := os.Create(latestFilePath)
		if err != nil {
			return fmt.Errorf("cannot create latest: %w", err)
		}
	}


	return nil
}

func WriteLog(level, message string, data interface{}) error {
	// Rotate if needed
	if err := rotateLogIfNeeded(); err != nil {
		return err
	}

	// Open latest log to append
	f, err := os.OpenFile(latestFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("write open error: %w", err)
	}
	defer f.Close()

	entry := LogEntry{
		Level: level,
		Message: message,
		Time: time.Now(),
		Data: data,
	}

	j, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	_, err = f.Write(append(j, '\n'))
	return err
}
