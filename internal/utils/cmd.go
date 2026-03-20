package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
)

// Check if there are command-line arguments
// If the first argument is "generate-api-key",
// call the GenerateAPIKey function from the key package to generate a new API key and save it to the database
func CommandTool(ctx context.Context, db sqlc.Querier) bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "generate-api-key":
			if err := GenerateAPIKey(ctx, db); err != nil {
				fmt.Printf("Error generating API key: %v\n", err)
			}
			return true
		case "revoke-api-key":
			if err := RevokeAPIKey(ctx, db); err != nil {
				fmt.Printf("Error when revoke API key: %v\n", err)
			}
			return true
		}
	}

	return false
}
