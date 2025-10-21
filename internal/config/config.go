package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DB_DSN     string
	JWTSecret  string
	SMTPUser   string
	SMTPPass   string
	SMTPHost string
	SMTPPort string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error with .env")
	}

	cfg := &Config{
		Port:      getEnv("APP_PORT", "8080"),
		DB_DSN:    getEnv("DB_DSN", ""),
		JWTSecret: getEnv("JWT_SECRET", "secret_key"),
		SMTPUser:  getEnv("SMTP_USER", ""),
		SMTPPass:  getEnv("SMTP_PASS", ""),
		SMTPHost:  getEnv("SMTP_HOST", ""),
		SMTPPort:  getEnv("SMTP_PORT", ""),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
