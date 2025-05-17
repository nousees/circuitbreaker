package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MaxRetries       int
	FailureThreshold uint32
	CBInterval       time.Duration
	CBTimeout        time.Duration
	ClientTimeout    time.Duration
	RetryDelay       time.Duration
	MaxRequests      uint32
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using default settings: %v", err)
	}

	return &Config{
		MaxRetries:       getEnvInt("MAX_RETRIES", 2),
		FailureThreshold: getEnvUint32("FAILURE_THRESHOLD", 3),
		CBInterval:       getEnvDuration("INTERVAL", 60*time.Second),
		CBTimeout:        getEnvDuration("TIMEOUT", 10*time.Second),
		ClientTimeout:    getEnvDuration("CLIENT_TIMEOUT", 5*time.Second),
		RetryDelay:       getEnvDuration("RETRY_DELAY", 500*time.Millisecond),
		MaxRequests:      getEnvUint32("MAX_REQUESTS", 2),
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvUint32(key string, defaultValue uint32) uint32 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(intValue)
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
