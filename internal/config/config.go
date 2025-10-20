package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisAddr   string
	JWTSecret   string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecret"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		RedisAddr:   redisAddr,
		JWTSecret:   jwtSecret,
	}
}
