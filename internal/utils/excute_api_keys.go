package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
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

func generateAPIKey() string {
	return fmt.Sprintf(
		"%s_%s",
		uuid.New().String(),
		randomString(16),
	)
}

func HashAPIKey(key string) string {
	pepper := os.Getenv("API_KEY_PEPPER")
	sum := sha256.Sum256([]byte(key + pepper))
	return hex.EncodeToString(sum[:])
}

func GenerateAPIKey(ctx context.Context, db sqlc.Querier) error {
	plaintext := generateAPIKey()
	keyHash := HashAPIKey(plaintext)

	genCmd := flag.NewFlagSet("generate-api-key", flag.ExitOnError)

	clientType := genCmd.String("type", "web", "Type of API key to generate (e.g., 'admin', 'user')")

	if err := genCmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	if *clientType == "" {
		return fmt.Errorf("client type is required")
	}

	if *clientType == "web" {
		if err := WriteEnv("VITE_API_KEY", plaintext); err != nil {
			return err
		}
	} else {
		key := fmt.Sprintf("%s_API_KEY", strings.ToUpper(*clientType))
		if err := WriteEnv(key, plaintext); err != nil {
			return err
		}
	}

	if err := db.CreateAPIKey(ctx, keyHash); err != nil {
		return err
	}

	return nil
}

func RevokeAPIKey(ctx context.Context, db sqlc.Querier) error {
	revokeCmd := flag.NewFlagSet("revoke-api-key", flag.ExitOnError)

	keyID := revokeCmd.String("id", "", "ID of the API key to revoke")

	if err := revokeCmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	keyHash := HashAPIKey(*keyID)

	if *keyID == "" {
		if err := db.RevokeAllAPIKeys(ctx); err != nil {
			return err
		}

		log.Println("All API keys revoked successfully")
		return nil
	}

	if err := db.RevokeAPIKeyByKey(ctx, keyHash); err != nil {
		return err
	}

	log.Printf("API key (%s) revoked successfully", *keyID)

	return nil
}
