package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds configuration variables
type Config struct {
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	JWT_SECRET  string
}

var AppConfig *Config

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	JWT_SECRET := os.Getenv("JWT_SECRET")

	// Initialize config struct
	AppConfig = &Config{
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_NAME:     DB_NAME,
		JWT_SECRET:  JWT_SECRET,
	}

	log.Println("Successfully loaded the .env file")

	return nil
}
