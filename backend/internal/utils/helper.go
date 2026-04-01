package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	pathEnv = ".env"
)

// loadEnv is a helper function that loads environment variables from a .env file using the godotenv package
func LoadEnv() {
	err := godotenv.Load(pathEnv)
	if err != nil {
		log.Println("No .env file found")
		panic(err)
	}
}

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

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting environment variable to int:", err)
		} else {
			return val
		}
	}

	return defaultValue
}

func WriteEnv(key string, value string) error {
	// Append to .env
	envLine := fmt.Sprintf("\n%s=%s\n", key, value)

	f, err := os.OpenFile(pathEnv, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(envLine); err != nil {
		return err
	}
	log.Printf("%s=%s", key, value)

	return nil
}

func SaveKeyToEnv(clientType, key string) error {
	if clientType == "web" {
		return WriteEnv("VITE_API_KEY", key)
	}

	envKey := fmt.Sprintf("%s_API_KEY", strings.ToUpper(clientType))

	log.Println("Saving API key to environment variable:", envKey)

	return WriteEnv(envKey, key)
}

func CheckUUID(id string) bool {
	err := uuid.Validate(id)
	if err != nil {
		return false
	}

	idUuid, err := uuid.Parse(id)
	if err != nil {
		return false
	}

	if idUuid == uuid.Nil {
		return false
	}

	return true
}

func GetKeyRedisAndConvertToInt(ctx context.Context, key string, rdb *redis.Client) (int, error) {
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, fmt.Errorf("Failed to get key from Redis: %v", err)
	}

	resultNum, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert Redis value to int: %v", err)
	}

	return resultNum, nil
}
