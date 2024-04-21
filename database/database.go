package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"github.com/marvy-O/fastchat/config"
)

var (
	db          *sql.DB
	once        sync.Once
	dbInitError error
)

// Database struct to hold the connection pool
type Database struct{}

// NewDatabase initializes the PostgreSQL connection if not already initialized
func NewDatabase() error {
	once.Do(func() {
		dbInitError = initializeDatabase()
	})

	if dbInitError != nil {
		return dbInitError
	}

	return nil
}

// InitDatabase initializes the PostgreSQL connection
func InitDatabase() error {
	if err := NewDatabase(); err != nil {
		return err
	}

	return nil
}

// ExecuteQuery is a function to run SQL queries
func ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	if err := NewDatabase(); err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return rows, nil
}

func executeSQLFile(sqlFilePath string, db *sql.DB) error {
	// Read the SQL file
	sqlFile, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}

	// Split the SQL script into individual statements
	sqlStatements := strings.Split(string(sqlFile), ";")

	// Execute each SQL statement
	for _, statement := range sqlStatements {
		if strings.TrimSpace(statement) == "" {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func initializeDatabase() error {
	// Replace these values with your PostgreSQL connection details
	// dbinfo := "user=postgres password=password dbname=fastchat sslmode=disable"
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.AppConfig.DB_USER, config.AppConfig.DB_PASSWORD, config.AppConfig.DB_NAME)
	conn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Check the connection
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Successfully connected to the database!")

	db = conn

	if err := executeSQLFile("./init.sql", db); err != nil {
		log.Fatalf("Error executing SQL file: %v\n", err)
	}

	log.Println("SQL file executed successfully")

	return nil
}
