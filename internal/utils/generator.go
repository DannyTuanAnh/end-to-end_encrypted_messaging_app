package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func GenerateApiKey(prefix string) string {
	return fmt.Sprintf(
		"%s_%s_%s",
		prefix,
		uuid.New().String(),
		randomString(16),
	)
}

func HashApiKey(key string) string {
	pepper := os.Getenv("API_KEY_PEPPER")
	sum := sha256.Sum256([]byte(key + pepper))
	return hex.EncodeToString(sum[:])
}
