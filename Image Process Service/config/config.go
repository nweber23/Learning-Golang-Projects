package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             string
	JWTSecret        string
	StoragePath      string
	RateLimitPerHour int
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	rateLimitStr := os.Getenv("RATE_LIMIT_PER_HOUR")
	rateLimit := 100
	if rateLimitStr != "" {
		if val, err := strconv.Atoi(rateLimitStr); err == nil {
			rateLimit = val
		}
	}
	return &Config{
		Port:             port,
		JWTSecret:        jwtSecret,
		StoragePath:      storagePath,
		RateLimitPerHour: rateLimit,
	}
}
