package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "modernc.org/sqlite"
)

// DBConfig holds database connection configuration
type DBConfig struct {
	DBType   string // "sqlite", "mysql", "postgres"
	Host     string
	Port     string
	DBName   string
	Username string
	Password string
	FilePath string // For SQLite
	IsLocal  bool
}

// AskForDatabaseConfig prompts the user for database configuration
func AskForDatabaseConfig(defaultRoute string) (DBConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	config := DBConfig{}

	fmt.Println("=== Database Configuration ===")
	fmt.Println("1. Local SQLite database")
	fmt.Println("2. Cloud database (MySQL/PostgreSQL)")
	fmt.Print("Select an option [1/2]: ")

	option, err := reader.ReadString('\n')
	if err != nil {
		return config, fmt.Errorf("error reading input: %v", err)
	}

	option = strings.TrimSpace(strings.ToLower(option))

	switch option {
	case "", "1":
		// Local SQLite database
		config.DBType = "sqlite"
		config.IsLocal = true

		fmt.Printf("Do you want to use the default database path (%s)? [Y/n]: ", defaultRoute)
		response, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))

		if response == "" || response == "y" || response == "yes" {
			config.FilePath = defaultRoute
			return config, nil
		}

		fmt.Print("Please input path to SQLite database: ")
		customPath, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading custom path: %v", err)
		}

		customPath = strings.TrimSpace(customPath)
		if customPath == "" {
			log.Println("Empty path provided, using default instead")
			config.FilePath = defaultRoute
		} else {
			config.FilePath = customPath
		}

	case "2":
		// Cloud database
		config.IsLocal = false

		fmt.Println("Select database type:")
		fmt.Println("1. MySQL")
		fmt.Println("2. PostgreSQL")
		fmt.Print("Select an option [1/2]: ")

		dbTypeOption, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}

		dbTypeOption = strings.TrimSpace(strings.ToLower(dbTypeOption))

		switch dbTypeOption {
		case "", "1":
			config.DBType = "mysql"
		case "2":
			config.DBType = "postgres"
		default:
			config.DBType = "mysql" // Default to MySQL
		}

		// Get database connection details
		fmt.Print("Host [localhost]: ")
		host, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}
		host = strings.TrimSpace(host)
		if host == "" {
			host = "localhost"
		}
		config.Host = host

		fmt.Print("Port [" + getDefaultPort(config.DBType) + "]: ")
		port, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}
		port = strings.TrimSpace(port)
		if port == "" {
			port = getDefaultPort(config.DBType)
		}
		config.Port = port

		fmt.Print("Database name: ")
		dbName, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}
		config.DBName = strings.TrimSpace(dbName)

		fmt.Print("Username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}
		config.Username = strings.TrimSpace(username)

		fmt.Print("Password: ")
		password, err := reader.ReadString('\n')
		if err != nil {
			return config, fmt.Errorf("error reading input: %v", err)
		}
		config.Password = strings.TrimSpace(password)

	default:
		return config, fmt.Errorf("invalid option: %s", option)
	}

	return config, nil
}

// getDefaultPort returns the default port for the given database type
func getDefaultPort(dbType string) string {
	switch dbType {
	case "mysql":
		return "3306"
	case "postgres":
		return "5432"
	default:
		return ""
	}
}

// Legacy function for backward compatibility
func AskForRoute(defaultRoute string) (string, bool, error) {
	config, err := AskForDatabaseConfig(defaultRoute)
	if err != nil {
		return "", false, err
	}

	if config.IsLocal && config.DBType == "sqlite" {
		return config.FilePath, config.FilePath == defaultRoute, nil
	}

	// For cloud databases, return an empty string and false to indicate non-default route
	return "", false, nil
}

// InitDB initializes the database connection based on configuration
func InitDB(filepath string) (*sql.DB, error) {
	// For backward compatibility, try to get the config first
	config, err := GetSavedConfig()
	if err != nil || config.DBType == "" {
		// If no saved config or error, use the filepath as SQLite database
		db, err := sql.Open("sqlite", filepath)
		if err != nil {
			return nil, err
		}

		// Check if database is accessible
		if err = db.Ping(); err != nil {
			return nil, err
		}

		log.Println("Connected to SQLite database")
		return db, nil
	}

	// Use the saved configuration
	return InitDBWithConfig(config)
}

// InitDBWithConfig initializes the database connection with the given configuration
func InitDBWithConfig(config DBConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error

	switch config.DBType {
	case "sqlite":
		db, err = sql.Open("sqlite", config.FilePath)
		if err != nil {
			return nil, err
		}
		log.Println("Connected to SQLite database")

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Username, config.Password, config.Host, config.Port, config.DBName)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		log.Println("Connected to MySQL database")

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.Username, config.Password, config.DBName)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}
		log.Println("Connected to PostgreSQL database")

	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.DBType)
	}

	// Check if database is accessible
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// SaveConfig saves the database configuration to a file
func SaveConfig(config DBConfig) error {
	file, err := os.Create("db_config.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write configuration to file
	fmt.Fprintf(file, "DBType=%s\n", config.DBType)
	fmt.Fprintf(file, "Host=%s\n", config.Host)
	fmt.Fprintf(file, "Port=%s\n", config.Port)
	fmt.Fprintf(file, "DBName=%s\n", config.DBName)
	fmt.Fprintf(file, "Username=%s\n", config.Username)
	fmt.Fprintf(file, "Password=%s\n", config.Password)
	fmt.Fprintf(file, "FilePath=%s\n", config.FilePath)
	fmt.Fprintf(file, "IsLocal=%v\n", config.IsLocal)

	return nil
}

// GetSavedConfig retrieves the saved database configuration
func GetSavedConfig() (DBConfig, error) {
	config := DBConfig{}

	file, err := os.Open("db_config.txt")
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "DBType":
			config.DBType = value
		case "Host":
			config.Host = value
		case "Port":
			config.Port = value
		case "DBName":
			config.DBName = value
		case "Username":
			config.Username = value
		case "Password":
			config.Password = value
		case "FilePath":
			config.FilePath = value
		case "IsLocal":
			config.IsLocal = value == "true"
		}
	}

	if err := scanner.Err(); err != nil {
		return config, err
	}

	return config, nil
}
