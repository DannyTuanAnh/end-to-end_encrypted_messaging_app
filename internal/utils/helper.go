package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
