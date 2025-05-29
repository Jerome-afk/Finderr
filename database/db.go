package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)

func AskForRoute(defaultRoute string) (string, bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Do you want to use the default route (%s)? [Y/n]: ", defaultRoute)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", false, fmt.Errorf("error reading input: %v", err)
	}

	// Normalize the response
	response = strings.TrimSpace(strings.ToLower(response))

	// If user answered y or pressed entered
	if response == "" || response == "y" || response == "yes" {
		return defaultRoute, true, nil
	}

	// Ask for custom route
	fmt.Print("Please input route to database: ")
	customRoute, err := reader.ReadString('\n')
	if err != nil {
		return "", false, fmt.Errorf("error reading custom route: %v", err)
	}

	customRoute = strings.TrimSpace(customRoute)
	if customRoute == "" {
		log.Println("Empty route provided, using default instead")
		return defaultRoute, true, nil
	}

	return customRoute, false, nil
}

// InitDB initializes the database connection
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Check if database is accessible
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}
