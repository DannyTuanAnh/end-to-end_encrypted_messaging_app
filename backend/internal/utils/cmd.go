package utils

import (
	"flag"
	"fmt"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
)

func ParseGenerateAPIKeyArgs(args []string) (*models.GenerateAPIKeyArgs, error) {
	genCmd := flag.NewFlagSet("generate-api-key", flag.ExitOnError)

	clientType := genCmd.String("type", "web", "Type of API key to generate (e.g., 'admin', 'user')")

	if err := genCmd.Parse(args); err != nil {
		return nil, err
	}

	if *clientType == "" {
		return nil, fmt.Errorf("client type is required")
	}

	return &models.GenerateAPIKeyArgs{
		ClientType: *clientType,
	}, nil
}

func ParseRevokeAPIKeyArgs(args []string) (*models.RevokeAPIKeyArgs, error) {
	revokeCmd := flag.NewFlagSet("revoke-api-key", flag.ExitOnError)

	keyID := revokeCmd.String("id", "", "ID of the API key to revoke")

	if err := revokeCmd.Parse(args); err != nil {
		return nil, err
	}

	if *keyID == "" {
		return &models.RevokeAPIKeyArgs{
			RevokeAll: true,
		}, nil
	}

	return &models.RevokeAPIKeyArgs{
		KeyID: *keyID,
	}, nil
}
