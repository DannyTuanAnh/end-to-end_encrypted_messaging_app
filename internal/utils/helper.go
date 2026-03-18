package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetEnvTime(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting environment variable to int:", err)
		} else {
			return time.Duration(val)
		}
	}

	return time.Duration(defaultValue)
}
