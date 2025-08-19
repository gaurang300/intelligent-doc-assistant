package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Gemini API configuration
	GeminiAPIKey string

	// Server configuration
	ServerPort string

	// Redis configuration (optional)
	RedisHost string
	RedisPort string
}

var (
	config *Config
	once   sync.Once
)

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
	once.Do(func() {
		// Load .env file if it exists
		godotenv.Load()

		config = &Config{
			DBHost:       getEnvOrDefault("DB_HOST", "localhost"),
			DBPort:       getEnvOrDefault("DB_PORT", "5432"),
			DBUser:       getEnvOrDefault("DB_USER", "postgres"),
			DBPassword:   getEnvOrDefault("DB_PASSWORD", ""),
			DBName:       getEnvOrDefault("DB_NAME", "docassistant"),
			GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
			ServerPort:   getEnvOrDefault("SERVER_PORT", "8080"),
			RedisHost:    getEnvOrDefault("REDIS_HOST", "localhost"),
			RedisPort:    getEnvOrDefault("REDIS_PORT", "6379"),
		}
	})
	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
