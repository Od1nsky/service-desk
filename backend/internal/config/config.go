package config

import (
	"os"
	"time"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration

	ServerPort  string
	CORSOrigins string
}

func Load() *Config {
	accessTTL, _ := time.ParseDuration(getenv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getenv("JWT_REFRESH_TTL", "720h"))
	return &Config{
		DBHost:      getenv("DB_HOST", "localhost"),
		DBPort:      getenv("DB_PORT", "5432"),
		DBUser:      getenv("DB_USER", "postgres"),
		DBPassword:  getenv("DB_PASSWORD", "postgres"),
		DBName:      getenv("DB_NAME", "servicedesk"),
		DBSSLMode:   getenv("DB_SSLMODE", "disable"),
		JWTSecret:   getenv("JWT_SECRET", "change-me"),
		AccessTTL:   accessTTL,
		RefreshTTL:  refreshTTL,
		ServerPort:  getenv("SERVER_PORT", "8080"),
		CORSOrigins: getenv("CORS_ORIGINS", "http://localhost:5173"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
