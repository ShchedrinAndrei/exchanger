package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	ServerPort          string
	ApiKey              string
	FetchIntervalMinute int
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "password"),
		DBName:              getEnv("DB_NAME", "currency_db"),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		ApiKey:              getEnv("FASTFOREX_API_KEY", ""),
		FetchIntervalMinute: getEnvAsInt("FETCH_INTERVAL_MINUTE", 1),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if valInt, err := strconv.Atoi(valStr); err == nil {
			return valInt
		}
	}
	return fallback
}
