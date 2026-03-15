package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port   string
	DBPath string
}

// LoadConfig loads configuration from environment variables
// Loads .env file if it exists, then reads environment variables
// Returns Config with default values if environment variables are not set
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		Port:   getEnvOrDefault("PORT", "8080"),
		DBPath: getEnvOrDefault("DB_PATH", "uptime_monitor.db"),
	}

	return cfg, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
