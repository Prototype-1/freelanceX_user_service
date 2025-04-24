package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret   string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	RedisAddr   string
	Port        string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	AppConfig = Config{
		JWTSecret:   getEnv("JWT_SECRET", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "freelanceX_user_service"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		Port:        getEnv("PORT", "8000"),
	}
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
