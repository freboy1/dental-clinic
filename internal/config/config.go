package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DB_DSN       string
	JWTSecret    string
	SMTPUser     string
	SMTPPass     string
	SMTPHost     string
	SMTPPort     string
	OpenAIAPIKey string
	OpenAIModel  string
	ResendAPIKey string
	FrontendURL  string
}

func LoadConfig() *Config {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("Error with .env")
	// }
	_ = godotenv.Load()

	cfg := &Config{
		Port:         getEnv("APP_PORT", "8080"),
		DB_DSN:       getEnv("DB_DSN", ""),
		JWTSecret:    getEnv("JWT_SECRET", "secret_key"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPass:     getEnv("SMTP_PASS", ""),
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", ""),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		ResendAPIKey: getEnv("ResendAPIKey", ""),
		FrontendURL:  getEnv("FrontendURL", ""),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
